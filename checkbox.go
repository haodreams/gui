/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2024-08-01 17:58:02
 * @FilePath: \goui\ui\Libel.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Checkbox struct {
	Widget[Checkbox]

	widget.Bool
	text string
}

func NewCheckbox(win *Window, text string) *Checkbox {
	m := new(Checkbox)
	m.text = text
	m.Init(win, m)
	m.Inset = win.CheckboxInset
	return m
}

func (m *Checkbox) SetValue(b bool) *Checkbox {
	m.Bool.Value = b
	return m
}

func (m *Checkbox) Layout(gtx layout.Context) layout.Dimensions {
	m.CheckDimensions(&gtx)
	return m.Inset.Layout(gtx, material.CheckBox(m.Theme(), &m.Bool, m.text).Layout)
}
