package main

import (
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
)

const (
	AppID          = "wx8375ef3da3481b01"
	AppSecret      = "5e83b0c44757577a106d001c04bc9ebb"
	Token          = "weiphp1127"
	EncodingAESKey = "P7aNQnLEFj8jEvMCYpujG2FYzpRxxO3rpMwHYhWt68L"
)

// 使用memcache保存access_token，也可选择redis或自定义cache
var wc = wechat.NewWechat()
var memory = cache.NewMemory()
var cfg = &offConfig.Config{
	AppID:          AppID,
	AppSecret:      AppSecret,
	Token:          Token,
	EncodingAESKey: EncodingAESKey,
	Cache:          memory,
}
var officialAccount = wc.GetOfficialAccount(cfg)

func WeiXinHandler(c *gin.Context) {
	// 传入request和responseWriter
	server := officialAccount.GetServer(c.Request, c.Writer)
	//server.SkipValidate(true)
	// 设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		// 回复消息：演示回复用户发送的消息
		logrus.Info("xxxxxxxxxx:", msg.Content)
		s := DoGPTRequest(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(s)}
	})
	// 处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		logrus.Warn(err)
		return
	}
	// 发送回复的消息
	err = server.Send()
	if err != nil {
		logrus.Warn(err)
		return
	}
}
