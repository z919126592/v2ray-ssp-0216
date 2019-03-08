package v2ray_sspanel_v3_mod_Uim_plugin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/rico93/v2ray-sspanel-v3-mod_Uim-plugin/client"
	"github.com/rico93/v2ray-sspanel-v3-mod_Uim-plugin/config"
	"github.com/rico93/v2ray-sspanel-v3-mod_Uim-plugin/webapi"
	"google.golang.org/grpc/status"
	"os"
	"runtime"
	"time"
	"v2ray.com/core/common/errors"
)

func init() {
	go func() {
		err := run()
		if err != nil {
			fatal(err)
		}
	}()
}
func checkAuth(panelurl string, db *webapi.Webapi) (bool, error) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(panelurl))
	cipherStr := md5Ctx.Sum(nil)
	current_md5 := hex.EncodeToString(cipherStr)
	re, _ := db.CheckAuth("https://auth.rico93.com", map[string]interface{}{"md5": current_md5})
	if re != nil {
		if re.Token != "" {
			auth := config.AESDecodeStr(re.Token, config.Key)
			return current_md5 == auth, nil
		} else {
			return false, newError("Auth failed or can't connect to the server")
		}

	} else {
		return false, newError("Auth failed or can't connect to the server")
	}

}
func run() error {

	err := config.CommandLine.Parse(os.Args[1:])

	cfg, err := config.GetConfig()
	if err != nil || *config.Test || cfg == nil {
		return err
	}

	// wait v2ray
	time.Sleep(3 * time.Second)
	db := &webapi.Webapi{
		WebToken:   cfg.PanelKey,
		WebBaseURl: cfg.PanelUrl,
	}
	ok, err := checkAuth(cfg.PanelUrl, db)
	if ok {
		go func() {
			apiInbound := config.GetInboundConfigByTag(cfg.V2rayConfig.Api.Tag, cfg.V2rayConfig.InboundConfigs)
			gRPCAddr := fmt.Sprintf("%s:%d", apiInbound.ListenOn.String(), apiInbound.PortRange.From)
			gRPCConn, err := client.ConnectGRPC(gRPCAddr, 10*time.Second)
			if err != nil {
				if s, ok := status.FromError(err); ok {
					err = errors.New(s.Message())
				}
				fatal(fmt.Sprintf("connect to gRPC server \"%s\" err: ", gRPCAddr), err)
			}
			newErrorf("Connected gRPC server \"%s\" ", gRPCAddr).AtWarning().WriteToLog()

			p, err := NewPanel(gRPCConn, db, cfg)
			if err != nil {
				fatal("new panel error", err)
			}

			p.Start()
		}()
		// Explicitly triggering GC to remove garbage
		runtime.GC()

	} else {
		fatal(err)
	}

	return nil
}

func newErrorf(format string, a ...interface{}) *errors.Error {
	return newError(fmt.Sprintf(format, a...))
}

func newError(values ...interface{}) *errors.Error {
	values = append([]interface{}{"SSPanelPlugin: "}, values...)
	return errors.New(values...)
}

func fatal(values ...interface{}) {
	newError(values...).AtError().WriteToLog()
	// Wait log
	time.Sleep(1 * time.Second)
	os.Exit(-2)
}
