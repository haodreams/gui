/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:57:29
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 16:46:44
 * @FilePath: \goui\ui\Table.go
 * @Description:
 */
package gui

import (
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type Tabler interface {
	//获取标题
	GetTitle(i int) (title string)
	//获取单元格数据
	GetItemText(record any, row, col int) (text string)
	//获取列宽度
	GetColumnWitdh(i int) (width float32)
	//获取列属性
	GetColumn(i int) *Column
	GetRow(row int) any
	//获取列个数
	GetColumnCount() (count int)
	//获取行数
	Size() (size int)
}

type Table struct {
	Widget[Table]
	component.GridState

	headerBorder    *widget.Border
	cellBorder      *widget.Border
	headers         []*widget.Clickable
	cells           []*widget.Clickable                         //单元格单击事件
	cellCb          func(gtx layout.Context, row, col, num int) //num 单击了几次
	rowChan         chan int                                    //存放需要刷新的行
	rowChanSize     int
	lastRowIdx      int //用于去重
	HealderClickIdx int //哪列被点击了

	creatCellCallback func(gtx layout.Context, row, col int) D

	rowIdx int //左键的索引
	colIdx int

	rRowIdx  int
	rColIdx  int
	clickNum int

	Tabler

	menu *Menu
}

func NewTable(win *Window, table Tabler) *Table {
	m := new(Table)
	m.Tabler = table
	m.Init(win, m)
	m.MaxWidth = -1
	m.MaxHeight = -1
	m.rowIdx = -1
	m.colIdx = -1
	m.rRowIdx = -1
	m.rColIdx = -1
	m.lastRowIdx = -1
	m.headerBorder = &widget.Border{
		Color: color.NRGBA{A: 255},
		Width: unit.Dp(0.2),
	}
	m.cellBorder = &widget.Border{
		Color: color.NRGBA{A: 150},
		Width: unit.Dp(0),
	}
	m.MaxHeight = 0
	m.MaxWidth = 0
	return m
}

func (m *Table) SetWidth(width int) *Table {
	m.MaxWidth = width
	return m
}
func (m *Table) SetHeight(height int) *Table {
	m.MaxHeight = height
	return m
}

// 创建自定义单元格回调函数
func (m *Table) SetCreateCellCallback(f func(gtx layout.Context, row, col int) D) *Table {
	m.creatCellCallback = f
	return m
}

func (m *Table) SetRowChannel(ch chan int, size int) *Table {
	m.rowChan = ch
	m.rowChanSize = size
	return m
}

func (m *Table) SetCellClick(f func(gtx layout.Context, row, col, num int)) *Table {
	m.cellCb = f
	return m
}

func (m *Table) GetSelectedCell() (int, int) {
	return m.rowIdx, m.colIdx
}

func (m *Table) SetMenu(menu *Menu) *Table {
	m.menu = menu
	return m
}

func (m *Table) Layout(gtx layout.Context) D {
	if len(m.headers) != m.GetColumnCount() {
		m.headers = make([]*widget.Clickable, m.GetColumnCount())
	}

	if len(m.cells) != m.Size()*m.GetColumnCount() {
		m.cells = make([]*widget.Clickable, m.Size()*m.GetColumnCount())
	}
	m.CheckDimensions(&gtx)

	return component.Table(m.Theme(), &m.GridState).Layout(gtx, m.Size(), m.GetColumnCount(),
		func(axis layout.Axis, index, constraint int) int {
			used := float32(0)
			defCount := 0
			defWidth := 70
			n := m.GetColumnCount()
			for i := range n {
				if m.GetColumnWitdh(i) == 0 {
					defCount++
					continue
				}
				used += m.GetColumnWitdh(i)
			}

			if defCount > 0 {
				free := (constraint - int(used)) / defCount
				if free > defWidth {
					defWidth = free
				}
			}
			switch axis {
			case layout.Horizontal:
				width := m.GetColumnWitdh(index)
				if width == 0 {
					return gtx.Dp(unit.Dp(defWidth))
				}
				return gtx.Dp(unit.Dp(width))
			default:
				return gtx.Dp(unit.Dp(30))
			}
		},
		func(gtx C, col int) D {
			return m.headerBorder.Layout(gtx, func(gtx C) D {
				//return inset.Layout(gtx, func(gtx C) D {
				column := m.GetColumn(col)
				click := m.headers[col]
				if click == nil {
					click = new(widget.Clickable)
					m.headers[col] = click
				}
				if click.Clicked(gtx) {
					m.HealderClickIdx = col
					if column.cb != nil {
						column.cb(col)
					}
				}
				return material.Clickable(gtx, click, func(gtx C) D {
					return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
						flatBtnText := material.Body1(m.Theme(), m.GetTitle(col))
						flatBtnText.MaxLines = 1
						flatBtnText.Truncator = "."
						return layout.Center.Layout(gtx, flatBtnText.Layout)
					})
				})
				//})
			})
		},
		func(gtx C, row, col int) D { //渲染单元格
			if m.lastRowIdx != row {
				m.lastRowIdx = row
				if m.rowChan != nil {
					if len(m.rowChan) < m.rowChanSize {
						m.rowChan <- row
					}
				}
			}
			record := m.GetRow(row)
			txt := m.GetItemText(record, row, col)
			idx := row*m.GetColumnCount() + col
			cell := m.cells[idx]
			if cell == nil {
				cell = &widget.Clickable{}
				m.cells[idx] = cell
			}
			if cell.Clicked(gtx) {
				if m.rowIdx == row && m.colIdx == col {
					m.clickNum++
				} else {
					m.clickNum = 1
				}
				m.rowIdx = row
				m.colIdx = col
				//m.win.Log(m.rowIdx, m.colIdx, m.clickNum)
				if m.cellCb != nil {
					m.cellCb(gtx, m.rowIdx, m.colIdx, m.clickNum)
				}
			}

			for {
				evt, ok := gtx.Source.Event(pointer.Filter{
					Target: cell,
					Kinds:  pointer.Press,
				})
				if !ok {
					break
				}
				e, ok := evt.(pointer.Event)
				if !ok {
					continue
				}
				if e.Buttons == pointer.ButtonSecondary {
					m.rRowIdx = row
					m.rColIdx = col
					//m.win.Log("right click", row, col)
				}
			}

			return material.Clickable(gtx, cell, func(gtx C) D {
				if m.menu != nil {
					m.menu.Clicked(gtx)
					if m.creatCellCallback != nil {
						//创建自定义单元格
						return m.creatCellCallback(gtx, row, col)
					}
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
						flatBtnText := material.Body1(m.Theme(), txt)
						if m.rowIdx == row {
							if m.colIdx == col {
								flatBtnText.Color = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
							} else {
								flatBtnText.Color = color.NRGBA{R: 0, G: 0, B: 180, A: 255}
							}
						}
						flatBtnText.MaxLines = 1
						flatBtnText.Truncator = "."
						return layout.Center.Layout(gtx, flatBtnText.Layout)
					})
				}

				if m.creatCellCallback != nil {
					//创建自定义单元格
					return m.creatCellCallback(gtx, row, col)
				}

				return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
					cellText := material.Body1(m.Theme(), txt)
					if m.rowIdx == row {
						if m.colIdx == col {
							cellText.Color = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
						} else {
							cellText.Color = color.NRGBA{R: 0, G: 0, B: 180, A: 255}
						}
					}
					column := m.GetColumn(col)
					cellText.Alignment = text.Alignment(column.Alignment)
					cellText.MaxLines = 1
					cellText.Truncator = "."
					return cellText.Layout(gtx)
					//return layout.Center.Layout(gtx, cellText.Layout)
				})
			})
		},
	)

}
