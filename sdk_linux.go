//go:build linux

package sdk

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include "dhnetsdk.h"

void *handle;
extern void cbDisConnect(LLONG lLoginID, char *pchDVRIP, LONG nDVRPort, LDWORD dwUser);
extern void screenshotCallBack(LLONG lLoginID, BYTE *pBuf, UINT RevLen, UINT EncodeType, DWORD CmdSerial, LDWORD dwUser);

int loadSo(const char* soPath){
	handle = dlopen(soPath, RTLD_LAZY);
    if (!handle) {
        return -1;
    }
	return 1;
}

BOOL client_init(fDisConnect cbDisConnect, LDWORD dwUser){
	return ((BOOL (*)(fDisConnect, LDWORD))dlsym(handle, "CLIENT_Init"))(cbDisConnect, dwUser);
}

LLONG client_LoginWithHighLevelSecurity(NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY* pstInParam, NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY* pstOutParam){
	return ((LLONG (*)(NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY*, NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY*))dlsym(handle, "CLIENT_LoginWithHighLevelSecurity"))(pstInParam, pstOutParam);
}

void client_SetSnapRevCallBack(fSnapRev OnSnapRevMessage, LDWORD dwUser){
	return ((void (*)(fSnapRev, LDWORD))dlsym(handle, "CLIENT_SetSnapRevCallBack"))(OnSnapRevMessage, dwUser);
}

BOOL client_SnapPictureEx(LLONG lLoginID, SNAP_PARAMS *par, int *reserved){
	return ((BOOL (*)(LLONG, SNAP_PARAMS*, int*))dlsym(handle, "CLIENT_SnapPictureEx"))(lLoginID, par, reserved);
}

BOOL client_Logout(LLONG lLoginID){
	return ((BOOL (*)(LLONG))dlsym(handle, "CLIENT_Logout"))(lLoginID);
}
*/
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

var loginIdMap = &sync.Map{}

type SDK struct {
	screenShotCh chan *ScreenshotResult
}

func NewSDK() (*SDK, error) {
	result := C.loadSo(C.CString(dllPath))
	if int(result) != 1 {
		return nil, errors.New(fmt.Sprint("load so err: ", result))
	}

	sdk := &SDK{screenShotCh: make(chan *ScreenshotResult)}
	err := sdk.Init()
	if err != nil {
		return nil, fmt.Errorf("error in Init: %v", err)
	}
	return sdk, nil
}

func (sdk *SDK) Init() error {
	code := C.client_init(C.fDisConnect(C.cbDisConnect), C.LDWORD(0))
	if code == 1 {
		return nil
	} else {
		return fmt.Errorf("error in CLIENT_Init, code %d", code)
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
	loginIdMap.Store(loginId, sdk)
	return loginId, &outParam.StuDeviceInfo, nil
}

func (sdk *SDK) LoginWithHighLevelSecurity(stuInParam *NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY, stuOutParam *NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY) (int, error) {
	ret := C.client_LoginWithHighLevelSecurity(
		(*C.NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY)(unsafe.Pointer(stuInParam)),
		(*C.NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY)(unsafe.Pointer(stuOutParam)),
	)
	return int(ret), nil
}

func (sdk *SDK) ScreenShot(loginId int, channel int) ([]byte, error) {
	C.client_SetSnapRevCallBack(C.fSnapRev(C.screenshotCallBack), C.LDWORD(0))
	p := &SNAP_PARAMS{
		Channel: uint32(channel),
		Quality: 1,
		Mode:    0,
	}
	var reserved int32
	code := C.client_SnapPictureEx(C.LLONG(loginId), (*C.SNAP_PARAMS)(unsafe.Pointer(p)), (*C.int)(unsafe.Pointer(&reserved)))
	if code != 1 {
		return nil, fmt.Errorf("error in SnapPictureEx, code: %d", code)
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
	code := C.client_Logout(C.LLONG(loginId))
	loginIdMap.Delete(loginId)
	if code != 1 {
		return fmt.Errorf("error in Logout, code: %d", code)
	}
	return nil
}
