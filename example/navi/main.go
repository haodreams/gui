/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-21 00:00:46
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-21 18:24:53
 * @FilePath: \gui\example\demo1.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"gioui.org/app"
	"gioui.org/unit"
	"gitee.com/haodreams/gui"
)

func Init() *gui.Window {
	win := gui.NewWindow()
	win.Option(
		app.Title("navi test"),
		app.Size(unit.Dp(800), unit.Dp(600)),
	)

	titles := []string{"首页", "通知", "推送", "我"}
	conts := make([]gui.Contenter, len(titles))
	for i, title := range titles {
		conts[i] = gui.NewEdit(win, title)
	}

	navi := gui.NewNavibar(win, titles, conts)

	win.SetContent(navi)
	return win
}

func main() {
	gui.Run(Init)
}
