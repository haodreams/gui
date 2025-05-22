/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-27 17:58:55
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 18:22:22
 * @FilePath: \dataviewe:\go\src\gitee.com\haodreams\gui\widget.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var (
	//White         = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	//Black         = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	// LightGreen    = color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0xff}
	// LightRed      = color.NRGBA{R: 0xff, G: 0x73, B: 0x73, A: 0xff}
	// LightYellow   = color.NRGBA{R: 0xff, G: 0xe0, B: 0x73, A: 0xff}
	// LightBlue     = color.NRGBA{R: 0x45, G: 0x89, B: 0xf5, A: 0xff}
	// LightPurple   = color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0xff}
	Border        = color.NRGBA{R: 0x6c, G: 0x6f, B: 0x76, A: 0xff}
	BorderFocused = color.NRGBA{R: 0x45, G: 0x89, B: 0xf5, A: 0xff}
	//DropDownMenuBg = color.NRGBA{R: 0x2b, G: 0x2d, B: 0x31, A: 0xff}
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Widget[T any] struct {
	win       *Window
	t         *T
	MaxWidth  int
	MaxHeight int
	MinWidth  int
	MinHeight int
	Inset     layout.Inset

	th *material.Theme //私有的主题，优先使用私有主题， 没有找到则找默认主题
}

func (m *Widget[T]) Theme() *material.Theme {
	if m.th == nil {
		return m.win.Theme()
	}
	return m.th
}

func (m *Widget[T]) SetTheme(th *material.Theme) *T {
	m.th = th
	return m.t
}

func (m *Widget[T]) Windos() *Window {
	return m.win
}

func (m *Widget[T]) Init(win *Window, t *T) {
	m.win = win
	m.t = t
}

func (m *Widget[T]) UniformInset(v unit.Dp) *T {
	m.Inset.Top = v
	m.Inset.Right = v
	m.Inset.Bottom = v
	m.Inset.Left = v
	return m.t
}

func (m *Widget[T]) SetInsetTop(v unit.Dp) *T {
	m.Inset.Top = v
	return m.t
}

func (m *Widget[T]) SetInsetBottom(v unit.Dp) *T {
	m.Inset.Bottom = v
	return m.t
}

func (m *Widget[T]) SetInsetLeft(v unit.Dp) *T {
	m.Inset.Left = v
	return m.t
}
func (m *Widget[T]) SetInsetRight(v unit.Dp) *T {
	m.Inset.Right = v
	return m.t
}

func (m *Widget[T]) CheckDimensions(gtx *layout.Context) {
	if m.MaxWidth > 0 {
		gtx.Constraints.Max.X = gtx.Dp(unit.Dp(m.MaxWidth))
	}
	if m.MaxHeight > 0 {
		gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(m.MaxHeight))
	}
	if m.MinWidth > 0 {
		gtx.Constraints.Min.X = gtx.Dp(unit.Dp(m.MinWidth))
	}
	if m.MinHeight > 0 {
		gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(m.MinHeight))
	}
}
func (m *Widget[T]) Resize(width, height int) *T {
	m.SetWidth(width)
	m.SetHeight(height)
	return m.t
}

func (m *Widget[T]) SetWidth(width int) *T {
	m.MaxWidth = width
	m.MinWidth = width
	return m.t
}

func (m *Widget[T]) SetHeight(height int) *T {
	m.MaxHeight = height
	m.MinHeight = height
	return m.t
}

func (m *Widget[T]) SetMaxWidth(width int) *T {
	m.MaxWidth = width
	return m.t
}

func (m *Widget[T]) SetMaxHeight(height int) *T {
	m.MaxHeight = height
	return m.t
}
func (m *Widget[T]) SetMinWidth(width int) *T {
	m.MinWidth = width
	return m.t
}

func (m *Widget[T]) SetMinHeight(height int) *T {
	m.MinHeight = height
	return m.t
}

func (m *Widget[T]) Default() *T {
	m.MaxWidth = 110
	m.MaxHeight = 30
	return m.t
}
