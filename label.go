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
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Label struct {
	Widget[Label]
	th * material.Theme
	material.LabelStyle
}

func NewLabel(win *Window, text string) *Label {
	m := new(Label)
	m.win = win
	m.Init(win, m)
	m.Inset = win.LabelInset
	m.LabelStyle = material.Body1(m.Theme(), text)
	return m
}

func (m *Label) EnableSelect(b bool) *Label {
	if b {
		if m.State == nil {
			m.State = new(widget.Selectable)
		}
		return m
	}
	m.State = nil
	return m
}

func (m *Label) SetText(text string) *Label {
	m.LabelStyle.Text = text
	return m
}

func (m *Label) SetAlignment(alig text.Alignment) *Label {
	m.LabelStyle.Alignment = alig
	return m
}

func (m *Label) SetColor(color color.NRGBA) *Label {
	m.LabelStyle.Color = color
	return m
}

func (m *Label)SetCaret(start, end int)  *Label{
	if m.State != nil{
		m.State.SetCaret(start,end)
	}
	return m
}

func (m *Label) SetStyle(f func (th * material.Theme, txt string) material.LabelStyle) *Label {
	if m.th == nil{
		m.th = m.Theme()
	}
	style:= f(m.th,m.Text)
	style.Color = m.Color
	style.State = m.State
	m.LabelStyle = style
	return m
}

func (m *Label) Layout(gtx layout.Context) layout.Dimensions {
	m.CheckDimensions(&gtx)
	return m.Inset.Layout(gtx, m.LabelStyle.Layout)
}
