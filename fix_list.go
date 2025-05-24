package gui

import (
	"io"

	"gioui.org/layout"
	"gioui.org/widget"
)

// 定长列表的控件，当超过最大容量后，会自动覆盖前面的数据
// 支持倒序显示
type FixList struct {
	Widget[FixList]
	widget.List
	FixArray[Contenter]
	desc bool //是否倒序显示
}

func NewFixList(win *Window, size int) *FixList {
	m := new(FixList)
	m.win = win
	m.Init(win, m)
	m.SetSize(size)
	return m
}

func (m *FixList) SetDesc(desc bool) *FixList {
	m.desc = desc
	return m
}

func (m *FixList) AddItem(item Contenter) *FixList {
	m.Push(&item)
	return m
}

func (m *FixList) SetSize(size int) *FixList {
	m.FixArray.SetSize(size)
	return m
}

func (m *FixList) Layout(gtx layout.Context) layout.Dimensions {
	m.CheckDimensions(&gtx)
	return m.List.Layout(gtx, m.Size(), func(gtx layout.Context, index int) layout.Dimensions {
		if m.desc {
			item := m.GetDesc(index)
			if item == nil {
				return D{}
			}
			return (*item).Layout(gtx)
		}
		item := m.Get(index)
		if item == nil {
			return D{}
		}
		return (*item).Layout(gtx)
	})
}

// 固定大小的列表
type FixArray[G any] struct {
	array []*G
	pos   uint //当前队列位置
	cap   uint //最大容量
	size  uint //当前数据个数
}

/**
 * @description: 新建一个固定大小的列表
 * @return {*}
 */
func NewFixArray[G any](size int) *FixArray[G] {
	m := new(FixArray[G])
	m.SetSize(size)
	return m
}

func (m *FixArray[G]) Reset() {
	m.pos = 0
	m.size = 0
}

/**
 * @description: 初始化队列的最大值
 * @param {int} maxLines
 * @return {*}
 */
func (m *FixArray[G]) SetSize(maxLines int) {
	if maxLines < 2 {
		maxLines = 2
	}
	if maxLines > 1000000 {
		maxLines = 1000000
	}
	if maxLines == int(m.cap) {
		return
	}
	m.array = make([]*G, maxLines)
	m.cap = uint(maxLines)
	m.Reset()
}

/**
 * @description: 获取队列数据大小
 * @param {*}
 * @return {*}
 */
func (m *FixArray[G]) Size() int {
	return int(m.size)
}

/**
 * @description: 压入数据
 * @param {*EventLog} e
 * @return {*}
 */
func (m *FixArray[G]) Push(e *G) {
	m.array[m.pos] = e
	m.pos++
	m.size++
	m.pos = m.pos % m.cap
	if m.size > m.cap {
		m.size = m.cap
	}
}

func (m *FixArray[G]) Get(index int) *G {
	if index < 0 || index >= int(m.size) {
		return nil
	}
	pos := ((m.cap + m.pos - m.size) + uint(index)) % m.cap
	return m.array[pos]
}

// 倒叙获取数据
func (m *FixArray[G]) GetDesc(index int) *G {
	if index < 0 || index >= int(m.size) {
		return nil
	}
	pos := ((m.cap + m.pos) - uint(index) - 1) % m.cap
	return m.array[pos]
}

/**
 * @description: 压入数据
 * @param {*EventLog} e
 * @return {*}
 */
func (m *FixArray[G]) Pop() (e *G, err error) {
	if m.size == 0 {
		return nil, io.EOF
	}

	pos := (m.cap + m.pos - m.size) % m.cap
	e = m.array[pos]
	m.size--
	return
}


// 先进后出
func (m *FixArray[G]) Stack(begin, limit int) []*G {
	count := m.Size()
	if count == 0 {
		return nil
	}
	if limit == -1 {
		limit = count
	}
	if limit > count {
		limit = count
	}
	end := min(begin+limit, count)
	if end <= begin {
		return nil
	}
	limit = end - begin

	rows := make([]*G, limit)
	ii := 0

		for i := begin; i < end; i++ {
			pos := ((m.cap + m.pos) - uint(i) - 1) % m.cap
			p := m.array[pos]
			rows[ii] = p
			ii++
		}
	return rows[:ii]
}
