/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-21 00:00:46
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 16:45:08
 * @FilePath: \gui\example\demo1.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"fmt"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gitee.com/haodreams/gui"
)

func Init() *gui.Window {
	win := gui.NewWindow()
	win.Option(
		app.Title("Counter"),
		app.Size(unit.Dp(800), unit.Dp(600)),
	)

	header := gui.NewContainer(win)

	header.AddWidget(gui.NewLabel(win, "Hello:"))

	header.AddWidget(gui.NewEdit(win, "test"))
	header.AddWidget(gui.NewEdit(win, "test2"))
	header.Add(gui.NewSpace(8))

	button := gui.NewButton(win, "Count", func() { win.Log(time.Now()) })
	header.Add(gui.NewSpace(8))

	button2 := gui.NewButton(win, "test", func() { win.Log(time.Now()) })

	header.AddWidget(button)
	header.Add(gui.NewSpace(8))

	header.Add(layout.Rigid(button2.Layout))

	tabs := gui.NewTabs(win)
	for i := 1; i <= 5; i++ {
		edit := gui.NewEdit(win, fmt.Sprintf("TEST Tab %d", i))
		edit.SetSigleLine(false)

		tab := gui.NewTab(win, fmt.Sprintf("Tab %d", i+1))
		tab.Add(layout.Flexed(1, edit.Layout))
		tabs.AddTab(tab)
	}

	root := gui.NewContainer(win)
	root.SetAxis(layout.Vertical)

	root.AddWidget(header)
	//root.Add(layout.Flexed(1, edit.Layout))
	root.Add(layout.Flexed(1, tabs.Layout))

	win.SetContent(root)
	return win
}

func main() {
	gui.Run(Init)
}
