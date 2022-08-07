package service

import (
	"context"
	"time"

	dingtalk_robot "github.com/JetBlink/dingtalk-notify-go-sdk"
	"github.com/gogf/gf/v2/frame/g"
)

func DinkTalkSend(ctx context.Context, msg string) (err error) {
	// 获取配置文件中的dingTalk.enable
	cfgEnable, err := g.Cfg().Get(ctx, "dingTalk.enable")
	if err != nil {
		g.Log().Error(ctx, "get config dingTalk.enable error: ", err)
		return
	}
	// 转换成bool
	enable := cfgEnable.Bool()
	//  如果配置文件中的dingTalk.enable为false，则不发送
	if !enable {
		g.Log().Info(ctx, "dingTalk.enable is false, skip")
		return
	}
	// 获取配置文件中的dingTalk.accessToken
	cfgDingTalkAccessToken, err := g.Cfg().Get(ctx, "dingTalk.accessToken")
	if err != nil {
		g.Log().Error(ctx, "get config dingTalk.accessToken error: ", err)
		return
	}
	// 获取配置文件中的dingTalk.secret
	cfgDingTalkSecret, err := g.Cfg().Get(ctx, "dingTalk.secret")
	if err != nil {
		g.Log().Error(ctx, "get config dingTalk.secret error: ", err)
		return
	}

	// 初始化钉钉机器人
	robot := dingtalk_robot.NewRobot(cfgDingTalkAccessToken.String(), cfgDingTalkSecret.String())
	// 发送消息
	robot.SendTextMessage(msg, []string{}, false)
	// 休眠5秒,防止触发频率过高
	g.Log().Info(ctx, "sleep 5 seconds")
	time.Sleep(5 * time.Second)
	return nil
}
