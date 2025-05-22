/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 20:57:36
 * @LastTextFieldors: wangjun haodreams@163.com
 * @LastTextFieldTime: 2024-07-28 10:38:30
 * @FilePath: \goui\ui\TextField.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/x/component"
)

type TextField struct {
	Widget[TextField]
	hint string
	component.TextField
	onSubmit func(string)
}

func NewTextField(win *Window, hint string) *TextField {
	m := new(TextField)
	m.Init(win, m)
	m.Inset = win.TextFieldInset
	m.SingleLine = true
	m.hint = hint
	return m
}

func (m *TextField) SetMask(mask rune) *TextField {
	m.Mask = mask
	return m
}

func (m *TextField) SetOnSubmit(f func(string)) *TextField {
	m.Submit = true
	m.onSubmit = f
	return m
}

func (m *TextField) SetText(text string) *TextField {
	m.TextField.SetText(text)
	return m
}

func (m *TextField) SetSigleLine(singleLine bool) *TextField {
	m.SingleLine = singleLine
	return m
}

func (m *TextField) Layout(gtx layout.Context) D {
	if m.TextField.Submit {
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

	// if m.SingleLine {
	// 	gtx.Constraints.Min = image.Pt(m.MaxWidth, m.MaxHeight)
	// 	gtx.Constraints.Max = image.Pt(m.MaxWidth, m.MaxHeight)
	// } else {
	// 	if m.MaxWidth != 0 {
	// 		gtx.Constraints.Min.X = m.MaxWidth
	// 		gtx.Constraints.Max.X = m.MaxWidth
	// 	}
	// 	if m.MaxHeight != 0 {
	// 		gtx.Constraints.Min.Y = m.MaxHeight
	// 		gtx.Constraints.Max.Y = m.MaxHeight
	// 	}
	// }
	m.CheckDimensions(&gtx)

	return m.TextField.Layout(gtx, m.Theme(), m.hint)
	// borderWidth := float32(0.5)
	// borderColor := color.NRGBA{A: 200}
	// switch {
	// case gtx.Source.Focused(&m.TextField):
	// 	borderColor = m.Theme().Palette.ContrastBg
	// 	borderWidth = 2
	// }
	// border := widget.Border{Color: borderColor, CornerRadius: unit.Dp(2), Width: unit.Dp(borderWidth)}
	// return border.Layout(gtx, func(gtx C) D {

	// 	return layout.UniformInset(unit.Dp(5)).Layout(gtx, material.TextFieldor(m.Theme(), &m.TextFieldor, m.text).Layout)
	// })
}
