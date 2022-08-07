package service

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func Checkbyconfig(ctx context.Context) (err error) {
	g.Log().Debug(ctx, "check by config ...")

	// 获取配置文件中的warnLeftDays
	cfgWarnLeftDays, err := g.Cfg().Get(ctx, "warnLeftDays")
	if err != nil {
		g.Log().Error(ctx, "get config warnLeftDays error: ", err)
		return
	}
	// 转换成int
	warnLeftDays := cfgWarnLeftDays.Int64()

	// 获取配置文件中的urls
	cfgurls, err := g.Cfg().Get(ctx, "urls")
	if err != nil {
		g.Log().Error(ctx, "get config urls error: ", err)
		return
	}
	// 转换成数组
	urls := cfgurls.Strings()
	g.Log().Debug(ctx, "urls: ", urls)
	// 定义全部正常标识变量
	allOk := true
	// 循环urls
	for _, url := range urls {
		g.Log().Info(ctx, url, "is checking ...")
		certInfo, err := GetCertInfo(ctx, url)
		if err != nil {
			g.Log().Error(ctx, "failed to get cert info: ", err)
			now := gtime.Now().String()
			msg := now + "[WARN] " + url + " failed to get cert info: " + err.Error()
			DinkTalkSend(ctx, msg)
			// return err
			allOk = false
			continue
		}
		g.Log().Info(ctx, url, "certInfo: ", certInfo)
		// 如果配置文件中的warnLeftDays为0，则不发送
		// if warnLeftDays == 0 {
		// 	g.Log().Info(ctx, "warnLeftDays is 0, skip")
		// 	continue
		// }
		// 计算证书还有多少天过期(使用unix时间戳)
		leftdays := (certInfo.NotAfterUnix - gtime.Now().Unix()) / (60 * 60 * 24)
		g.Log().Info(ctx, url, "leftdays: ", leftdays)
		// 如果证书还有多少天过期小于配置文件中的warnLeftDays，则发送钉钉消息
		if leftdays < warnLeftDays {
			g.Log().Info(ctx, url, "will expire in ", leftdays, " days", ", send dingTalk message")
			now := gtime.Now().String()
			msg := now + "[WARN] ssl cert in " + url + " will expire in " + strconv.Itoa(int(leftdays)) + " days"
			DinkTalkSend(ctx, msg)
			allOk = false
		}
	}
	if allOk {
		g.Log().Info(ctx, "all ok")
		// 获取当前时间
		now := gtime.Now().String()
		msg := now + " all ssl cert is ok"
		DinkTalkSend(ctx, msg)
	}

	return nil
}
