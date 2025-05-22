/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 18:30:18
 * @FilePath: \goui\ui\Libel.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type Progress struct {
	Widget[Progress]

	process float32
}

func NewProgress(win *Window) *Progress {
	m := new(Progress)
	m.Init(win, m)
	m.Inset = win.SelectInset
	return m
}

func (m *Progress) SetProgress(progress float32) *Progress {
	m.process = progress
	return m
}
func (m *Progress) Layout(gtx layout.Context) D {
	m.CheckDimensions(&gtx)
	return m.Inset.Layout(gtx, material.ProgressBar(m.Theme(), m.process).Layout)
	//return m.Inset.Layout(gtx, material.Progress(m.Theme(), &m.Bool, m.text).Layout)
}
