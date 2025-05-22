package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Split struct {
	Color  color.NRGBA
	Height unit.Dp
	Left   int
}

func NewSplit(c color.NRGBA) *Split {
	return &Split{
		Color:  c,
		Height: 1,
	}
}

func (m *Split) Layout(gtx C) D {
	width := gtx.Constraints.Max.X
	d := image.Point{X: width, Y: gtx.Dp(m.Height)}

	defer clip.UniformRRect(image.Rectangle{Min: image.Pt(m.Left, 0), Max: image.Pt(width, d.Y)}, 0).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: d}
}
