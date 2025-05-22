package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func NewSpace(dp ...int) layout.FlexChild {
	if len(dp) == 0 {
		return layout.Flexed(1, func(gtx C) D {
			return layout.Spacer{}.Layout(gtx) //D{} // 占位符，用于填充剩余空间
		})
	}
	width := 0
	height := 0
	switch len(dp) {
	case 1:
		height = dp[0]
	case 2:
		width, height = dp[0], dp[1]
	}
	return layout.Rigid(layout.Spacer{Width: unit.Dp(width), Height: unit.Dp(height)}.Layout)
}

type space struct {
	layout.Spacer
}

func (m *space) Layout(gtx layout.Context) D {
	return m.Spacer.Layout(gtx) //D{} // 占位符，用于填充剩余空间
}

func Space() *space {
	return new(space)
}

// func NewSpace() layout.FlexChild {
// 	return layout.Flexed(1, func(gtx C) D {
// 		return D{} // 占位符，用于填充剩余空间
// 	})
// }
