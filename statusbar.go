/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-21 18:14:57
 * @FilePath: \goui\ui\Libel.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type statusBar struct {
	win  *Window
	text string
	show bool
}

func (m *statusBar) Init(win *Window) *statusBar {
	m.win = win
	return m
}

func (m *statusBar) SetStatusBarText(text string) {
	m.text = text
}

func (m *statusBar) ShowStatusBar(b bool) {
	m.show = b
}

func (m *statusBar) Layout(gtx layout.Context) D {
	return material.Body1(m.win.theme, m.text).Layout(gtx)
}
