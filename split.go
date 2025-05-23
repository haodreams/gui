package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var (
	//默认的分割线颜色
	defaultSplitColor = color.NRGBA{R: 200, G: 200, B: 200, A: 200} //灰色
)

type Split struct {
	color  color.NRGBA
	height unit.Dp
	Left   int
}

func NewSplit() *Split {
	return &Split{
		color:  defaultSplitColor,
		height: 1,
	}
}

func (m *Split) SetHeight(h unit.Dp) *Split {
	m.height = h
	return m
}

func (m *Split) SetColor(c color.NRGBA) *Split {
	m.color = c
	return m
}

func (m *Split) Layout(gtx C) D {
	width := gtx.Constraints.Max.X
	d := image.Point{X: width, Y: gtx.Dp(m.height)}

	defer clip.UniformRRect(image.Rectangle{Min: image.Pt(m.Left, 0), Max: image.Pt(width, d.Y)}, 0).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: m.color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: d}
}
