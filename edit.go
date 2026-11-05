/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 20:57:36
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-27 12:17:32
 * @FilePath: \goui\ui\edit.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"image"
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
	blur      func(string) //失去光标的回调函数
	focused   bool         //当前是否有焦点，判断失去焦点的条件
	textSize  unit.Sp      //字体大小
	color     *color.NRGBA //字体颜色
	onSubmit  func(string)
	clearable bool          //是否显示清空按钮
	btnClear  widget.Clickable //清空按钮
}

func NewEdit(win *Window, hint string) *Edit {
	m := new(Edit)
	m.Init(win, m)
	m.Inset = win.EditInset
	m.SingleLine = true
	m.hint = hint
	return m
}

// SetClearButton 设置是否显示清空按钮(X), 有内容时点击清空输入框
func (m *Edit) SetClearButton(clearable bool) *Edit {
	m.clearable = clearable
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

func (m *Edit)SetColor(c *color.NRGBA) *Edit {
	m.color = c
	return m
}

func (m *Edit)SetTextSize(size unit.Sp) *Edit {
	m.textSize = size
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
		//return m.Inset.Layout(gtx, material.Editor(m.Theme(), &m.Editor, m.hint).Layout)
		es := material.Editor(m.Theme(), &m.Editor, m.hint)
		if m.textSize > 0 {
			es.TextSize = m.textSize
		}
		if m.color != nil {
			es.Color = *m.color
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return m.Inset.Layout(gtx, es.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				return m.layoutClear(gtx)
			}),
		)
	})
}

// layoutClear 绘制清空按钮, 输入框有内容且开启了清空按钮时显示
func (m *Edit) layoutClear(gtx C) D {
	if !m.clearable || m.Editor.Text() == "" {
		return D{}
	}
	if m.btnClear.Clicked(gtx) {
		m.Editor.SetText("")
		m.Windos().Invalidate()
	}
	size := gtx.Dp(unit.Dp(18))
	gtx.Constraints.Min = image.Pt(size, size)
	gtx.Constraints.Max = image.Pt(size, size)
	return material.Clickable(gtx, &m.btnClear, func(gtx C) D {
		gtx.Constraints.Min = image.Pt(size, size)
		gtx.Constraints.Max = image.Pt(size, size)
		return CloseIcon.Layout(gtx, m.Theme().Fg)
	})
}
