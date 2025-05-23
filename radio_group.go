/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2024-07-28 11:16:12
 * @FilePath: \goui\ui\Libel.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type RadioGroup struct {
	Widget[RadioGroup]
	layout.Flex
	widget.Enum
	// keys   []string
	// values []string
	radios []layout.FlexChild
}

func NewRadioGroup(win *Window, keys, values []string) *RadioGroup {
	m := new(RadioGroup)
	m.win = win
	m.radios = make([]layout.FlexChild, len(keys))
	for i := range keys {
		m.radios[i] = layout.Rigid(material.RadioButton(m.Theme(), &m.Enum, keys[i], values[i]).Layout)
	}
	m.Init(win, m)
	return m
}

func (m *RadioGroup) SetValue(v string) *RadioGroup {
	m.Value = v
	return m
}

func (m *RadioGroup) SetAxis(axis layout.Axis) *RadioGroup {
	m.Axis = axis
	return m
}

func (m *RadioGroup) Layout(gtx layout.Context) D {
	return m.Flex.Layout(gtx, m.radios...)
}
