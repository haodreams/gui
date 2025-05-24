/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-21 00:00:46
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 11:30:07
 * @FilePath: \gui\example\demo1.go
 * @Description:
 */
package main

import (
	"embed"
	"image"
	"log"

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

	boss.man = NewPageMain(win, "main", "首页")
	boss.msg = NewPageMsg(win, "msg", "通知")
	boss.push = NewPagePush(win, "push", "推送")
	boss.me = NewPageMe(win, "me", "我的")
}

func (m *Boss) GetPages() (titles []string, contents []gui.Contenter) {
	titles = []string{boss.man.Title(), boss.msg.Title(), boss.push.Title(), boss.me.Title()}
	contents = []gui.Contenter{boss.man, boss.msg, boss.push, boss.me}
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
