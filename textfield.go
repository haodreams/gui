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
	blur     func(string) //失去光标的回调函数
	focused  bool         //当前是否有焦点，判断失去焦点的条件
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

func (m *TextField) SetBlur(f func(string)) *TextField {
	m.blur = f
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
	if m.blur != nil {
		if gtx.Source.Focused(&m.Editor) {
			m.focused = true
		} else {
			if m.focused {
				m.focused = false
				m.blur(m.TextField.Text())
			}
		}
	}

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

	m.CheckDimensions(&gtx)
	return m.TextField.Layout(gtx, m.Theme(), m.hint)
}
