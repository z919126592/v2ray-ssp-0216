package db

import (
	"github.com/imroc/req"
	"github.com/rico93/v2ray-sspanel-v3-mod_Uim-plugin/model"
	"github.com/rico93/v2ray-sspanel-v3-mod_Uim-plugin/speedtest"
)

type NodeinfoResponse struct {
	Ret  uint            `json:"ret"`
	Data *model.NodeInfo `json:"data"`
}
type PostResponse struct {
	Ret  uint   `json:"ret"`
	Data string `json:"data"`
}
type UsersResponse struct {
	Ret  uint              `json:"ret"`
	Data []model.UserModel `json:"data"`
}
type AllUsers struct {
	Ret  uint
	Data map[string]model.UserModel
}

type DisNodenfoResponse struct {
	Ret  uint                 `json:"ret"`
	Data []*model.DisNodeInfo `json:"data"`
}
type AuthResponse struct {
	Ret   uint   `json:"ret"`
	Token string `json:"token"`
}

var id2string = map[uint]string{
	0: "server_address",
	1: "port",
	2: "alterid",
	3: "protocol",
	4: "protocol_param",
	5: "path",
	6: "host",
	7: "inside_port",
	8: "server",
}
var maps = map[string]interface{}{
	"server_address": "",
	"port":           "",
	"alterid":        "16",
	"protocol":       "tcp",
	"protocol_param": "",
	"path":           "",
	"host":           "",
	"inside_port":    "",
	"server":         "",
}

type Db interface {
	GetNodeInfo(uint) (*NodeinfoResponse, error)

	GetDisNodeInfo(uint) (*DisNodenfoResponse, error)

	GetALLUsers(*model.NodeInfo) (*AllUsers, error)
	Post(string, map[string]interface{}, map[string]interface{}) (*req.Resp, error)
	UploadSystemLoad(uint) bool
	UpLoadUserTraffics(uint, []model.UserTrafficLog) bool
	UploadSpeedTest(uint, []speedtest.Speedresult) bool
	UpLoadOnlineIps(uint, []model.UserOnLineIP) bool

	CheckAuth(string, map[string]interface{}) (*AuthResponse, error)
}
