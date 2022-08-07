package cmd

import (
	"context"
	"hjm-certcheck/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gcron"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start cert check server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			buildInfo := gbuild.Info()
			g.Log().Info(ctx, "start hjm-certcheck server ...")
			g.Log().Info(ctx, "HJM-CERTCHECK,https://gitee.com/hjmcloud/hjm-certcheck")
			g.Log().Info(ctx, "version: ", buildInfo.Git)
			cronPattern, err := g.Cfg().Get(ctx, "cron.pattern")
			if err != nil {
				g.Log().Error(ctx, "get config cron.pattern error: ", err)
				return
			}
			// cronPattern 不能为空
			if cronPattern.IsEmpty() {
				g.Log().Error(ctx, "config cron.pattern is empty")
				return
			}
			cronRunOnStart, err := g.Cfg().Get(ctx, "cron.runOnStart")
			if err != nil {
				g.Log().Error(ctx, "get config cron.runOnStart error: ", err)
				return
			}
			if cronRunOnStart.Bool() {
				g.Log().Info(ctx, "cron.runOnStart is true, run cron job on start")
				err = service.Checkbyconfig(ctx)
				if err != nil {
					g.Log().Error(ctx, "run cron job on start error: ", err)
				}
			}
			g.Log().Info(ctx, "cron.pattern: ", cronPattern.String())
			// 增加定时任务 certCheck
			g.Log().Info(ctx, "add cron job certCheck ...")
			_, err = gcron.Add(ctx, cronPattern.String(), func(ctx context.Context) {
				err = service.Checkbyconfig(ctx)
				if err != nil {
					g.Log().Error(ctx, "add cron job error: ", err)
					return
				}
			}, "certCheck")
			g.Log().Info(ctx, "add cron job certCheck success")

			// service.Checkbyconfig(ctx)

			// 防退出
			g.Log().Info(ctx, "hjm-certcheck server is running ...")
			select {}
		},
	}
)
