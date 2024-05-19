//go:build linux

package sdk

/*
#include "dhnetsdk.h"
extern void cbDisConnect(LLONG lLoginID, char *pchDVRIP, LONG nDVRPort, LDWORD dwUser);
extern void screenshotCallBack(LLONG lLoginID, BYTE *pBuf, UINT RevLen, UINT EncodeType, DWORD CmdSerial, LDWORD dwUser);
*/
import "C"
import (
	"time"
	"unsafe"
)

//export cbDisConnect
func cbDisConnect(lLoginID C.LLONG, pchDVRIP *C.char, nDVRPort C.LONG, dwUser C.LDWORD) {}

//export screenshotCallBack
func screenshotCallBack(loginId C.LLONG, buf *C.BYTE, revLen C.UINT, encodeType C.UINT, cmdSerial C.DWORD, dwUser C.LDWORD) {
	sdkInterface, has := loginIdMap.Load(int(loginId))
	if !has {
		return
	}
	sdk := sdkInterface.(*SDK)
	picBuf := (*[1 << 30]byte)(unsafe.Pointer(buf))[:int(revLen):int(revLen)]
	select {
	case sdk.screenShotCh <- &ScreenshotResult{picBuf, int(loginId)}:
		return
	case <-time.After(time.Second * 1):
		return
	}
}
