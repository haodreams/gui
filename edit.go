/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 20:57:36
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2024-07-28 18:13:14
 * @FilePath: \goui\ui\edit.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Edit struct {
	Widget[Edit]
	hint string
	widget.Editor
	blur    func(string) //失去光标的回调函数
	focused bool         //当前是否有焦点，判断失去焦点的条件

	onSubmit func(string)
}

func NewEdit(win *Window, hint string) *Edit {
	m := new(Edit)
	m.Init(win, m)
	m.Inset = win.EditInset
	m.SingleLine = true
	m.hint = hint
	return m
}

func (m *Edit) SetBlur(f func(string)) *Edit {
	m.blur = f
	return m
}

func (m *Edit) SetMask(mask rune) *Edit {
	m.Mask = mask
	return m
}

func (m *Edit) SetOnSubmit(f func(string)) *Edit {
	m.Submit = true
	m.onSubmit = f
	return m
}

func (m *Edit) SetText(text string) *Edit {
	m.Editor.SetText(text)
	return m
}

func (m *Edit) SetSigleLine(singleLine bool) *Edit {
	m.SingleLine = singleLine
	return m
}

func (m *Edit) Layout(gtx layout.Context) layout.Dimensions {
	if m.blur != nil {
		if gtx.Source.Focused(&m.Editor) {
			m.focused = true
		} else {
			if m.focused {
				m.focused = false
				m.blur(m.Text())
			}
		}
	}

	if m.Editor.Submit {
		for {
			e, ok := m.Editor.Update(gtx)
			if !ok {
				break
			}
			if e, ok := e.(widget.SubmitEvent); ok {
				if m.onSubmit != nil {
					m.onSubmit(e.Text)
				}
			}
		}
	}

	m.CheckDimensions(&gtx)

	borderWidth := float32(0.5)
	borderColor := color.NRGBA{A: 200}
	switch {
	case gtx.Source.Focused(&m.Editor):
		borderColor = m.Theme().Palette.ContrastBg
		borderWidth = 2
	}
	border := widget.Border{Color: borderColor, CornerRadius: unit.Dp(2), Width: unit.Dp(borderWidth)}
	return border.Layout(gtx, func(gtx C) D {
		return m.Inset.Layout(gtx, material.Editor(m.Theme(), &m.Editor, m.hint).Layout)
	})
}
