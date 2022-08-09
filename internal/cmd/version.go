package cmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gutil"
)

type sVersion struct {
	Name     string //程序名称
	Homepage string //程序主页
	Version  string //程序版本
	GoFrame  string //goframe version
	Golang   string //golang version
	Git      string //git commit id
	Time     string //build datetime
}

var (
	Version = gcmd.Command{
		Name:  "version",
		Usage: "version",
		Brief: "show version",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			info := gbuild.Info()
			// gutil.Dump(info)
			data := gbuild.Data()
			// gutil.Dump(data)

			// 生成sVersion结构体
			res := sVersion{
				Name:     "hjm-certcheck",
				Homepage: "https://hjmcloud.cn/hjm-certcheck",
				Version:  data["version"].(string),
				GoFrame:  info.GoFrame,
				Golang:   info.Golang,
				Git:      info.Git,
				Time:     info.Time,
			}
			gutil.Dump(res)
			return nil
		},
	}
)

// init 初始化模块
func init() {
	Main.AddCommand(&Version)
}
