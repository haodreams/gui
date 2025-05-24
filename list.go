/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 19:02:22
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

func (m *List) AddItem(item Contenter) *List {
	m.items = append(m.items, item)
	return m
}

func (m *List) Layout(gtx layout.Context) layout.Dimensions {
	m.CheckDimensions(&gtx)
	return m.List.Layout(gtx, len(m.items), func(gtx layout.Context, index int) layout.Dimensions {
		return m.items[index].Layout(gtx)
	})
}

