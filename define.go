package sdk

import "unsafe"

var loginErrorMap = map[int]string{
	1:  "账号或密码错误",
	2:  "用户名不存在",
	3:  "登录超时",
	4:  "重复登录",
	5:  "帐号被锁定",
	6:  "帐号被列入禁止名单",
	7:  "系统忙,资源不足",
	8:  "子连接失败",
	9:  "主连接失败",
	10: "超过最大连接数",
	11: "只支持3代协议",
	12: "设备未插入U盾或U盾信息错误",
	13: "客户端IP地址没有登录权限",
	18: "设备账号未初始化，无法登陆",
}

// 定义数据类型
type C_DWORD uint32
type C_ENUM int32
type C_LLONG int64

// NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY 结构体
type NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY struct {
	DwSize     C_DWORD
	SzIP       [64]byte
	NPort      int32
	SzUserName [64]byte
	SzPassword [64]byte
	EmSpecCap  int32
	ByReserved [4]byte
	PCapParam  unsafe.Pointer
	EmTLSCap   C_ENUM
}

// NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY 结构体
type NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY struct {
	DwSize        C_DWORD
	StuDeviceInfo NET_DEVICEINFO_Ex
	NError        int32
	ByReserved    [132]byte
}

// NET_DEVICEINFO_Ex 结构体
type NET_DEVICEINFO_Ex struct {
	SSerialNumber    [48]byte
	NAlarmInPortNum  int32
	NAlarmOutPortNum int32
	NDiskNum         int32
	NDVRType         int32
	NChanNum         int32
	ByLimitLoginTime byte
	ByLeftLogTimes   byte
	BReserved        [2]byte
	NLockLeftTime    int32
	Reserved         [24]byte
}

type SNAP_PARAMS struct {
	Channel   uint32
	Quality   uint32
	ImageSize uint32
	Mode      uint32
	InterSnap uint32
	CmdSerial uint32
	Reserved  [4]uint32
}
