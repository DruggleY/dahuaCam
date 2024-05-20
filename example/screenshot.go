package main

import (
	"encoding/json"
	"fmt"
	sdk "github.com/DruggleY/dahuaCam"
	"os"
)

func main() {
	sdk, err := sdk.NewSDK()
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
