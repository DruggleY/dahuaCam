package sdk

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
	"time"
	"unsafe"
)

//go:embed dhnetsdk.dll
var dllFile []byte
var dllPath string

func init() {
	file, err := os.CreateTemp("", "dhnetsdk.dll.*")
	if err != nil {
		panic(err)
	}
	_, err = file.Write(dllFile)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	dllPath = file.Name()
}

type SDK struct {
	dll          *syscall.DLL
	screenShotCh chan *ScreenshotResult
}

func NewSDK() (*SDK, error) {
	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		return nil, fmt.Errorf("error in LoadDLL: %v", err)
	}
	sdk := &SDK{dll: dll, screenShotCh: make(chan *ScreenshotResult)}
	err = sdk.Init()
	if err != nil {
		return nil, fmt.Errorf("error in Init: %v", err)
	}
	return sdk, nil
}

func (sdk *SDK) Proc(name string) *syscall.Proc {
	proc, err := sdk.dll.FindProc(name)
	if err != nil {
		panic(err)
	}
	return proc
}

func (sdk *SDK) Init() error {
	cbDisConnect := func(lLoginID int64, pchDVRIP *byte, nDVRPort int64, dwUser int64) uintptr {
		return 0
	}
	code, _, err := sdk.Proc("CLIENT_Init").Call(windows.NewCallback(cbDisConnect), uintptr(0))
	if code == 1 {
		return nil
	} else {
		return fmt.Errorf("error in CLIENT_Init, code %d: %v", code, err)
	}
}

func (sdk *SDK) Login(ip string, port int, username string, password string) (int, *NET_DEVICEINFO_Ex, error) {
	inParam := &NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY{
		DwSize: C_DWORD(224),
		NPort:  int32(port),
	}
	copy(inParam.SzIP[:], ip)
	copy(inParam.SzUserName[:], username)
	copy(inParam.SzPassword[:], password)
	outParam := &NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY{
		DwSize: C_DWORD(unsafe.Sizeof(NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY{})),
	}
	loginId, err := sdk.LoginWithHighLevelSecurity(inParam, outParam)
	if err != nil {
		return 0, nil, fmt.Errorf("error in call login: %v", err)
	}
	if loginId == 0 {
		errorMsg, has := loginErrorMap[int(outParam.NError)]
		if !has {
			errorMsg = "未知错误"
		}
		return 0, nil, fmt.Errorf("%s", errorMsg)
	}
	return loginId, &outParam.StuDeviceInfo, nil
}

// LoginWithHighLevelSecurity 函数
func (sdk *SDK) LoginWithHighLevelSecurity(stuInParam *NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY, stuOutParam *NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY) (int, error) {
	ret, _, err := sdk.Proc("CLIENT_LoginWithHighLevelSecurity").Call(
		uintptr(unsafe.Pointer(stuInParam)),
		uintptr(unsafe.Pointer(stuOutParam)),
	)
	if errors.Is(err, syscall.Errno(0)) {
		err = nil
	}
	return int(ret), err
}

type ScreenshotResult struct {
	pic     []byte
	loginId int
}

func (sdk *SDK) screenshotCallBack(loginId int64, buf *byte, revLen, encodeType, cmdSerial, dwUser int64) uintptr {
	if loginId == 0 {
		return 0
	}
	picBuf := (*[1 << 30]byte)(unsafe.Pointer(buf))[:revLen:revLen]
	select {
	case sdk.screenShotCh <- &ScreenshotResult{picBuf, int(loginId)}:
		return 0
	case <-time.After(time.Second * 1):
		return 0
	}
}

func (sdk *SDK) ScreenShot(loginId int, channel int) ([]byte, error) {
	_, _, err := sdk.Proc("CLIENT_SetSnapRevCallBack").Call(windows.NewCallback(sdk.screenshotCallBack), uintptr(0))
	if !errors.Is(err, syscall.Errno(0)) {
		return nil, fmt.Errorf("error in SetSnapRevCallBack: %v", err)
	}
	p := &SNAP_PARAMS{
		Channel: uint32(channel),
		Quality: 1,
		Mode:    0,
	}
	var reserved int32
	_, _, err = sdk.Proc("CLIENT_SnapPictureEx").Call(uintptr(loginId), uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(&reserved)))
	if !errors.Is(err, syscall.Errno(0)) {
		return nil, fmt.Errorf("error in SnapPictureEx: %v", err)
	}

	timeLimit := time.After(time.Second * 10)
	for {
		select {
		case result := <-sdk.screenShotCh:
			if result.loginId == loginId {
				return result.pic, nil
			}
		case <-timeLimit:
			return nil, errors.New("timeout")
		}
	}
}

func (sdk *SDK) Logout(loginId int) error {
	ret, _, err := sdk.Proc("CLIENT_Logout").Call()
	if !errors.Is(err, syscall.Errno(0)) {
		return fmt.Errorf("error in logout: %v", err)
	}
	if ret == 0 {
		return syscall.GetLastError()
	}
	return nil
}

func main() {
	sdk, err := NewSDK()
	if err != nil {
		panic(err)
	}

	// 调用高安全级别登录函数
	loginId, info, err := sdk.Login("", 37777, "admin", "admin123")
	if err != nil {
		fmt.Println("Login failed:", err)
		return
	}

	j, _ := json.Marshal(info)
	fmt.Println(string(j))

	b, err := sdk.ScreenShot(loginId, 0)
	if err != nil {
		fmt.Println("ScreenShot failed:", err)
	}
	err = os.WriteFile("screenshot.jpg", b, 0644)
	if err != nil {
		fmt.Println("WriteFile failed:", err)
	}

	if err := sdk.Logout(loginId); err != nil {
		fmt.Println("Logout failed:", err)
	}
}
