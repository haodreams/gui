/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2024-08-01 17:41:15
 * @FilePath: \goui\ui\Libel.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type Label struct {
	Widget[Label]
	material.LabelStyle
	text string
}

func NewLabel(win *Window, text string) *Label {
	m := new(Label)
	m.text = text
	m.win = win
	m.Init(win, m)
	m.Inset = win.LabelInset
	m.LabelStyle = material.Body1(m.Theme(), m.text)
	return m
}

func (m *Label) SetAlignment(alig text.Alignment) *Label {
	m.LabelStyle.Alignment = alig
	return m
}

func (m *Label) SetH6() *Label {
	m.LabelStyle = material.H6(m.Theme(), m.text)
	return m
}

func (m *Label) Layout(gtx layout.Context) layout.Dimensions {
	m.CheckDimensions(&gtx)
	return m.Inset.Layout(gtx, m.LabelStyle.Layout)
}
