/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-26 09:54:57
 * @FilePath: \goui\ui\Libel.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type List struct {
	Widget[List]
	widget.List
	items []Contenter
}

func NewList(win *Window) *List {
	m := new(List)
	m.win = win
	m.Init(win, m)
	return m
}

func (m *List) SetAxis(axis layout.Axis) *List {
	m.Axis = axis
	return m
}

func (m *List) AddItem(item Contenter) *List {
	m.items = append(m.items, item)
	return m
}

func (m *List) Layout(gtx C) D {
	m.CheckDimensions(&gtx)
	return m.List.Layout(gtx, len(m.items), func(gtx C, index int) D {
		if m.items[index] == nil{
			return D{}
		}
		return m.items[index].Layout(gtx)
	})
}
