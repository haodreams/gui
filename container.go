/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:58:38
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 12:49:34
 * @FilePath: \goui\gui\container.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
)

type Container struct {
	Widget[Container]
	layout.Flex
	Title     string
	childrens []layout.FlexChild
}

func NewContainer(win *Window) *Container {
	m := new(Container)
	m.Init(win, m)
	return m
}

func (m *Container) SetTitle(title string) *Container {
	m.Title = title
	return m
}

func (m *Container) SetAlignment(align layout.Alignment) *Container {
	m.Flex.Alignment = align
	return m
}

func (m *Container) SetSpacing(spacing layout.Spacing) *Container {
	m.Flex.Spacing = spacing
	return m
}

func (m *Container) SetHorizontal() *Container {
	m.Flex.Axis = layout.Horizontal
	return m
}

func (m *Container) SetVertical() *Container {
	m.Flex.Axis = layout.Vertical
	return m
}

func (m *Container) SetAxis(axis layout.Axis) *Container {
	m.Flex.Axis = axis
	return m
}

func (m *Container) Layout(gtx layout.Context) layout.Dimensions {
	m.CheckDimensions(&gtx)
	return m.Flex.Layout(gtx, m.childrens...)
}

func (m *Container) Add(child layout.FlexChild) *Container {
	m.childrens = append(m.childrens, child)
	return m
}

func (m *Container) AddWidget(widget Contenter, weight ...float32) *Container {
	var child layout.FlexChild
	if len(weight) > 0 {
		child = layout.Flexed(weight[0], widget.Layout)
	} else {
		child = layout.Rigid(widget.Layout)
	}

	m.Add(child)
	return m
}
