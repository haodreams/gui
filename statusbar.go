/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-08-01 09:49:35
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
	text1  material.LabelStyle
	text2 material.LabelStyle
	layout.Flex
	show bool
	only1 bool
}

func (m *statusBar) Init(win *Window) *statusBar {
	m.win = win
	m.text1 = material.Body1(m.win.theme, "")
	m.text2 = material.Body1(m.win.theme, "")
	m.Flex.Axis =layout.Horizontal
	return m
}

func (m *statusBar)SetStatusBarOnly1(b bool){
	m.only1 = b
}

func (m *statusBar) SetStatusBarText(text string, pos...int) {
	if m.only1{
		m.text1.Text = text
	}
	if len(pos) > 0{
		if pos[0] == 0{
			m.text1.Text = text
			return
		}
		m.text2.Text = text
		return
	}
	m.text1.Text = text
}

func (m *statusBar) ShowStatusBar(b bool) {
	m.show = b
}

func (m *statusBar) Layout(gtx layout.Context) D {
	if m.only1{
		return 	m.text1.Layout(gtx)

	}
	return m.Flex.Layout(gtx, layout.Flexed(0.66, m.text1.Layout), layout.Flexed(0.34, m.text2.Layout))
}
