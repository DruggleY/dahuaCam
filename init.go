package sdk

import (
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed bin/win64/dhnetsdk.dll
var win64DLL []byte

//go:embed bin/win32/dhnetsdk.dll
var win32DLL []byte

//go:embed bin/linux64/libdhnetsdk.so
var linux64So []byte

//go:embed bin/linux32/libdhnetsdk.so
var linux32So []byte

var dllPath string

type ScreenshotResult struct {
	pic     []byte
	loginId int
}

func init() {
	sdkMap := map[string][]byte{
		"windows_amd64": win64DLL,
		"windows_386":   win32DLL,
		"linux_amd64":   linux64So,
		"linux_386":     linux32So,
	}
	filename := map[string]string{
		"windows": "dhnetsdk.dll",
		"linux":   "libdhnetsdk.so",
	}[runtime.GOOS]
	tmpDir := os.TempDir()
	soPath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(soPath, sdkMap[runtime.GOOS+"_"+runtime.GOARCH], 0755)
	if err != nil {
		panic(err)
	}
	dllPath = soPath
}
