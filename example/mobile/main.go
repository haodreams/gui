/*
- @Author: wangjun haodreams@163.com
- @Date: 2024-07-28 16:59:19

  - @LastEditors: wangjun haodreams@163.com

  - @LastEditTime: 2024-08-03 14:22:33

- @FilePath: \gui\example\mobile\main.go
- @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
*/
package main

import (
	"io"
	"os"

	"gioui.org/app"
	_ "gioui.org/app/permission/networkstate"
	_ "gioui.org/app/permission/storage"
	"gioui.org/widget"
	"gitee.com/haodreams/gui"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func Init() *gui.Window {
	//sock.CreateUdpListener(":39000")
	w := gui.NewWindow()
	w.Option(
		app.Title("新锐实时数据库"),
	)
	w.Theme().TextSize = 14

	head := gui.NewContainer(w)
	cont := gui.NewContainer(w).SetVertical()
	root := gui.NewContainer(w).SetVertical()
	root.AddWidget(head)
	root.AddWidget(cont, 1)

	ic, _ := widget.NewIcon(icons.ContentAdd)

	head.AddWidget(gui.NewSelect(w,
		gui.NewSelectOption("数据库").WithValue("remote"),
		gui.NewSelectOption("本地数据").WithValue("local")).SetWidth(150))
	head.AddWidget(gui.NewLabel(w, "测试1"))
	head.AddWidget(gui.NewEdit(w, "测试2"), 1)
	head.AddWidget(gui.NewButton(w, "测试3", func() {
		w.MessageBox("test", "hello world!", gui.ModalError)
	}).SetIcon(ic).SetVertical().SetInsetTop(0).SetInsetBottom(0).SetTextSize(10))

	edit1 := gui.NewEdit(w, "测试输入")
	edit2 := gui.NewEdit(w, "测试输入2")

	cont.AddWidget(edit1)
	cont.AddWidget(edit2)
	cont.AddWidget(gui.NewTextField(w, "测试hint"))
	cont.AddWidget(gui.NewCheckbox(w, "测试复选"))
	path := "/storage/emulated/0/Android/media/com.gitee.mobile/files/"
	err := os.MkdirAll(path, 0755)
	if err != nil {
		edit1.SetText(err.Error())
	} else {
		edit1.SetText("创建目录成功")
	}
	cont.AddWidget(gui.NewButton(w, "读文件", func() {
		data, err := os.ReadFile(path + "test.txt")
		if err != nil {
			edit1.SetText(err.Error())
			return
		}
		edit1.SetText(string(data))
		w.MessageBox("test", "hello world1!", gui.ModalError).SetHeight(-1)
	}))

	cont.AddWidget(gui.NewButton(w, "写文件", func() {
		err := os.WriteFile(path+"test.txt", []byte("hello world"), 0644)
		if err != nil {
			edit2.SetText(err.Error())
		} else {
			edit2.SetText("write file success.")
		}
	}))
	cont.AddWidget(gui.NewButton(w, "测试5", func() {
		w.MessageBox("test", "hello world!", gui.ModalError)
	}))
	cont.AddWidget(gui.NewButton(w, "测试6", func() {
		go func() {
			file, err := w.ChooseFile("txt", "conf")
			w.Log("---------------------333--------", err)
			if err != nil {
				w.MessageBox("test", "hello world!", gui.ModalError)

				return
			}
			data, err := io.ReadAll(file)
			w.Log("-----------------------------", err, string(data))
			file.Close()
			w.MessageBox("aaa", string(data), gui.ModalInfo)
		}()
	}))
	cont.AddWidget(gui.NewRadioGroup(w, []string{"选项1", "选项2", "选项3"}, []string{"选项1", "选项2", "选项3"}))

	w.SetContent(root)
	w.ShowStatusBar(true)
	return w

}

// gogio -target android -icon logo.png -signkey "E:\android\test.keystore" -signpass 19861029 . && adb install mobile.apk
func main() {
	gui.Run(Init)
}
