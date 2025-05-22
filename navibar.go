/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:58:38
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 16:46:30
 * @FilePath: \goui\gui\Navibar.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
)

type navibar struct {
	Widget[navibar]
	barFlex   layout.Flex
	flex      layout.Flex
	btns      []*Button
	selected  int
	contents  []Contenter
	childrens []layout.FlexChild
	bar       layout.FlexChild
	space     layout.FlexChild
}

func NewNavibar(win *Window, titles []string, contents []Contenter) *navibar {
	m := new(navibar)
	m.Init(win, m)
	m.barFlex.Spacing = layout.SpaceBetween
	m.flex.Axis = layout.Vertical

	last := len(titles)
	m.btns = make([]*Button, last)
	m.contents = make([]Contenter, last)
	for i, title := range titles {
		btn := NewButton(m.win, title, func() { m.do(i) }).Resize(50, 30).UniformInset(0)
		btn.SetBackground(m.win.theme.Bg)
		btn.SetColor(m.win.theme.Fg)
		m.btns[i] = btn
		m.contents[i] = contents[i]
		if i == 0 {
			cont := NewContainer(m.win).Add(NewSpace(20, 20)).AddWidget(btn)
			m.childrens = append(m.childrens, layout.Rigid(cont.Layout))
			continue
		}
		if i == last {
			cont := NewContainer(m.win).AddWidget(btn).Add(NewSpace(20, 20))
			m.childrens = append(m.childrens, layout.Rigid(cont.Layout))
			continue
		}
		m.childrens = append(m.childrens, layout.Rigid(btn.Layout))
	}
	m.bar = layout.Rigid(func(gtx C) D {
		return m.barFlex.Layout(gtx, m.childrens...)
	})
	m.space = NewSpace(5, 5)
	return m
}

func (m *navibar) SetSelected(i int) *navibar {
	m.selected = i
	return m
}

func (m *navibar) Layout(gtx layout.Context) D {
	return m.flex.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return m.contents[m.selected].Layout(gtx)
		}),
		m.bar,
		m.space,
	)
}

func (m *navibar) do(i int) {
	m.selected = i
	m.win.Log("switch to ", m.selected)
}
