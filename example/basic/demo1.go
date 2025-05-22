/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-21 00:00:46
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 14:34:40
 * @FilePath: \gui\example\demo1.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"log"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gitee.com/haodreams/gui"
)

func Init() *gui.Window {
	win := gui.NewWindow(
		gui.WithLog(log.Println),
	)
	win.Option(
		app.Title("Counter"),
		app.Size(unit.Dp(800), unit.Dp(600)),
	)

	progress := gui.NewProgress(win)

	header := gui.NewContainer(win)

	header.AddWidget(gui.NewLabel(win, "Hello:"))

	header.AddWidget(gui.NewEdit(win, "").SetWidth(100))
	header.AddWidget(gui.NewEdit(win, "test2").SetWidth(120))
	header.Add(gui.NewSpace(8))

	button := gui.NewButton(win, "60%", func() {
		win.Log(time.Now())
		progress.SetProgress(0.6) // 60%
		win.SetStatusBarText("OK")

	}).SetIcon(gui.FileFolderOpen)

	header.Add(gui.NewSpace(8))

	header.AddWidget(button)
	header.Add(gui.NewSpace(8))

	mode := 0
	header.AddWidget(gui.NewButton(win, "mode", func() {
		mode++
	}))

	edit := gui.NewEdit(win, "test2")
	edit.SetSigleLine(false)

	dr := gui.NewSelect(win,
		gui.NewSelectOption("test22").WithValue("test222"),
		gui.NewSelectOption("test22").WithValue("test222"),
		gui.NewSelectOption("test22").WithValue("test222"),
		gui.NewSelectOption("test223").WithValue("test2232"),
	)

	header.AddWidget(dr)

	root := gui.NewContainer(win)
	root.SetAxis(layout.Vertical)

	root.AddWidget(header)

	checkbox := gui.NewCheckbox(win, "test")
	root.AddWidget(checkbox)
	header.Add(gui.NewSpace(8, 10))

	root.AddWidget(progress)

	radioGroup1 := gui.NewRadioGroup(win, []string{"key1", "key2"}, []string{"value1", "value2"}).SetAxis(layout.Vertical)
	root.AddWidget(radioGroup1)
	radioGroup2 := gui.NewRadioGroup(win, []string{"key1", "key2"}, []string{"value3", "value4"})
	root.AddWidget(radioGroup2)

	root.AddWidget(gui.NewButton(win, "info", func() {
		win.Info("test", "hello world!")
	}))
	root.AddWidget(gui.NewButton(win, "warn", func() {
		win.Warn("test", "hello world!")
	}))
	root.AddWidget(gui.NewButton(win, "error", func() {
		win.Error("test", "hello world!")
	}))

	root.AddWidget(gui.NewButton(win, "prompt", func() {
		win.Prompt("test", "hello world!", func(text string) {
			win.Log(text)
		})
	}))

	win.ShowStatusBar(true)
	win.SetStatusBarText("ready22222222222222222222.")

	// statusbar := gui.NewStatusBar(win, 2)
	// statusbar.SetText(0, "ready.")
	// statusbar.SetText(1, "ping.")
	root.Add(layout.Flexed(1, edit.Layout))
	// root.AddWidget(statusbar)
	win.SetContent(root)
	return win
}

func main() {
	gui.Run(Init)
}
