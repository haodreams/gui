/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-21 00:00:46
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 22:55:06
 * @FilePath: \gui\example\demo1.go
 * @Description:
 */
package main

import (
	"embed"
	"fmt"
	"image"
	"log"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/unit"
	"gitee.com/haodreams/gui"
)

// 嵌入图片资源

//go:embed images/*
var res embed.FS

var boss = new(Boss)

// gogio -target android -icon logo.png -signkey "E:\android\test.keystore" -signpass xxxxx . && adb install alarm.apk

type Boss struct {
	//pages
	win  *gui.Window
	man  *PageMain
	msg  *PageMsg
	push *PagePush
	me   *PageMe

	//icons
	info  image.Image
	warn  image.Image
	erron image.Image
}

func Init() *gui.Window {
	win := gui.NewWindow(gui.WithLog(log.Println))
	win.Option(
		app.Title("Alarm test"),
		app.Size(unit.Dp(400), unit.Dp(700)),
	)

	boss.Init(win)

	titles, contents := boss.GetPages()
	navi := gui.NewNavibar(win, titles, contents)
	navi.SetSelected(1)
	win.SetContent(navi)
	return win
}

func main() {
	gui.Run(Init)
}

func (m *Boss) Init(win *gui.Window) {
	m.win = win
	m.win.SetOnClose(func() {
		log.Println("on close")
	})
	log.Println("data dir:", win.DataDir())
	m.loadJpeg()
	m.win.Shield.Show()
	m.win.Shield.SetContent(NewLoginDialog(win))

	m.man = NewPageMain(win, "main", "首页")
	m.msg = NewPageMsg(win, "msg", "通知")
	m.push = NewPagePush(win, "push", "推送")
	m.me = NewPageMe(win, "me", "我的")

	now := time.Now().Unix()
	m.PushMsg("", now-86400*5, "info", "Info", "这是一条提示信息")
	m.PushMsg("", now-86500, "warn", "Warn", "这是一条警告信息")
	m.PushMsg("", now, "error", "Error", "这是一条错误信息")

}

func (m *Boss) PushMsg(app string, time int64, level, title, text string) {
	do := func(card *gui.Card) {
		log.Println("cart id:", card.ID(), card.Msg())
		m.win.Info(card.Title(), card.Msg())
	}

	if app != "" {
		title = fmt.Sprintf("%s:%s", app, title)
	}
	cart := gui.NewCard(m.win, time, title, text).SetImage(m.info).SetOnClick(do)
	cart.SetLevel(level)
	level = strings.ToLower(level)
	switch level {
	case "info":
		cart.SetImage(m.info)
	case "warn":
		cart.SetImage(m.warn)
	case "error":
		cart.SetImage(m.erron)
	default: //TODO 加载自定义的图像
		cart.SetImage(m.info)
	}
	m.msg.list.AddItem(cart)
}

func (m *Boss) GetPages() (titles []string, contents []gui.Contenter) {
	titles = []string{m.man.Title(), m.msg.Title(), m.push.Title(), m.me.Title()}
	contents = []gui.Contenter{m.man, m.msg, m.push, m.me}
	return
}

func (m *Boss) loadJpeg() {
	data, err := res.ReadFile("images/info.jpg")
	if err != nil {
		log.Println(err)
	}
	m.info, _ = gui.LoadJpeg(data)

	data, err = res.ReadFile("images/warn.jpg")
	if err != nil {
		log.Println(err)
	}
	m.warn, _ = gui.LoadJpeg(data)

	data, err = res.ReadFile("images/error.jpg")
	if err != nil {
		log.Println(err)
	}
	m.erron, _ = gui.LoadJpeg(data)
}
