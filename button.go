/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2024-08-03 13:45:19
 * @FilePath: \goui\ui\button.go
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

type Button struct {
	Widget[Button]
	color    *color.NRGBA
	bgColor  *color.NRGBA
	icon     *widget.Icon
	iconRect bool
	spacer   int
	Axis     layout.Axis
	widget.Clickable
	f        func()
	text     string
	textSize unit.Sp
}

func NewButton(win *Window, text string, f func()) *Button {
	m := new(Button)
	m.text = text
	m.f = f
	m.Init(win, m)
	m.Inset = win.ButtonInset
	return m
}

func (m *Button) SetHorizontal() *Button {
	m.Axis = layout.Horizontal
	return m
}

func (m *Button) SetTextSize(textSize unit.Sp) *Button {
	m.textSize = textSize
	return m
}

func (m *Button) SetVertical() *Button {
	m.Axis = layout.Vertical
	return m
}
func (m *Button) SetSpacer(spacer int) *Button {
	m.spacer = spacer
	return m
}

func (m *Button) SetRectIcon(iconRect bool) *Button {
	m.iconRect = iconRect
	return m
}

func (m *Button) SetIcon(icon *widget.Icon) *Button {
	m.icon = icon
	return m
}

func (m *Button) SetBackground(bg color.NRGBA) *Button {
	m.bgColor = &bg
	return m
}

func (m *Button) SetColor(c color.NRGBA) *Button {
	m.color = &c
	return m
}

func (m *Button) Layout(gtx layout.Context) layout.Dimensions {
	if m.f != nil {
		if m.Clicked(gtx) {
			m.f()
		}
	}

	m.CheckDimensions(&gtx)
	if m.icon == nil {
		btn := material.Button(m.Theme(), &m.Clickable, m.text)
		if m.bgColor != nil {
			btn.Background = *m.bgColor
		}
		if m.color != nil {
			btn.Color = *m.color
		}
		btn.Inset = m.Inset
		return btn.Layout(gtx)
	}

	if m.text == "" {
		if m.iconRect {
			return material.Clickable(gtx, &m.Clickable, func(gtx C) D {
				if m.color != nil {
					return m.Inset.Layout(gtx, func(gtx C) D {
						return m.icon.Layout(gtx, *m.color)
					})
				}
				return m.Inset.Layout(gtx, func(gtx C) D {
					return m.icon.Layout(gtx, m.Theme().Fg)
				})
			})
		}
		btn := material.IconButton(m.Theme(), &m.Clickable, m.icon, m.text)
		if m.bgColor != nil {
			btn.Background = *m.bgColor
		}
		if m.color != nil {
			btn.Color = *m.color
		}
		btn.Inset = layout.UniformInset(2)
		return btn.Layout(gtx)
	}
	//图标和文字都有两个都画
	return material.ButtonLayout(m.win.theme, &m.Clickable).Layout(gtx, func(gtx C) D {
		top := m.Inset.Top - 2
		bottom := m.Inset.Bottom - 2
		if top < 0 {
			top = 0
		}
		if bottom < 0 {
			bottom = 0
		}
		return layout.Inset{Top: top, Bottom: bottom, Left: m.Inset.Left, Right: m.Inset.Right}.Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: m.Axis, Alignment: layout.Middle}
			layIcon := layout.Rigid(func(gtx C) D {
				var d D
				if m.color != nil {
					d = m.icon.Layout(gtx, *m.color)
				} else {
					d = m.icon.Layout(gtx, m.win.theme.ContrastFg)
				}

				if m.Axis == layout.Horizontal {
					return layout.Inset{Right: unit.Dp(m.spacer)}.Layout(gtx, func(gtx C) D {
						return d
					})
				}
				return layout.Inset{Bottom: unit.Dp(m.spacer)}.Layout(gtx, func(gtx C) D {
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				l := material.Body1(m.Theme(), m.text)
				l.Color = m.win.theme.Palette.ContrastFg
				if m.color != nil {
					l.Color = *m.color
				}
				if m.Axis == layout.Horizontal {
					return layout.Inset{Left: unit.Dp(m.spacer)}.Layout(gtx, func(gtx C) D {
						return l.Layout(gtx)
					})
				}
				if m.textSize > 0 {
					l.TextSize = m.textSize
				}
				return layout.Inset{Top: unit.Dp(m.spacer)}.Layout(gtx, func(gtx C) D {
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}
