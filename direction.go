package gui

import "gioui.org/layout"

type Direction struct {
	layout.Widget
	layout.Direction
}

func NewDirection(d layout.Direction, widget layout.Widget) *Direction {
	m := new(Direction)
	m.Direction = d
	m.Widget = widget
	return m
}

func (m *Direction) Layout(gtx layout.Context) layout.Dimensions {
	return m.Direction.Layout(gtx, m.Widget)
}
