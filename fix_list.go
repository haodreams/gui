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
	FixSlice[Contenter]
	desc bool //是否倒序显示
}

func NewFixList(win *Window, size int) *FixList {
	m := new(FixList)
	m.win = win
	m.Init(win, m)
	m.SetSize(size)
	return m
}

func (m *FixList) SetAxis(axis layout.Axis) *FixList {
	m.Axis = axis
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
	m.FixSlice.SetSize(size)
	return m
}

func (m *FixList) Layout(gtx layout.Context) layout.Dimensions {
	m.CheckDimensions(&gtx)
	return m.List.Layout(gtx, m.Len(), func(gtx layout.Context, index int) layout.Dimensions {
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

// 固定大小的队列
type FixSlice[G any] struct {
	array []*G
	head  uint32 // 队首位置
	tail  uint32 // 队尾的下一个位置
	cap   uint32 // 最大容量
	size  uint32 // 当前数据个数
}

// 创建新的固定大小队列
func NewFixSlice[G any](size int) *FixSlice[G] {
	m := &FixSlice[G]{}
	m.SetSize(size)
	return m
}

// 重置队列
func (m *FixSlice[G]) Reset() {
	m.head = 0
	m.tail = 0
	m.size = 0
	// 帮助GC回收内存
	for i := range m.array {
		m.array[i] = nil
	}
}

// 设置队列大小
func (m *FixSlice[G]) SetSize(maxLines int) {
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
	m.cap = uint32(maxLines)
	m.Reset()
}

// 获取队列大小
func (m *FixSlice[G]) Len() int {
	return int(m.size)
}

// 检查队列是否为空
func (m *FixSlice[G]) IsEmpty() bool {
	return m.size == 0
}

// 检查队列是否已满
func (m *FixSlice[G]) IsFull() bool {
	return m.size == m.cap
}

// 入队操作
func (m *FixSlice[G]) Push(e *G) {
	if m.IsFull() {
		// 队列已满，覆盖队首元素
		m.array[m.head] = e
		m.head = (m.head + 1) % m.cap
	} else {
		m.array[m.tail] = e
		m.tail = (m.tail + 1) % m.cap
		m.size++
	}
}

// 获取指定位置的元素
func (m *FixSlice[G]) Get(index int) *G {
	if index < 0 || index >= int(m.size) {
		return nil
	}
	pos := (m.head + uint32(index)) % m.cap
	return m.array[pos]
}

// 倒序获取指定位置的元素
func (m *FixSlice[G]) GetDesc(index int) *G {
	if index < 0 || index >= int(m.size) {
		return nil
	}
	pos := (m.tail - uint32(index) - 1 + m.cap) % m.cap
	return m.array[pos]
}

// 出队操作
func (m *FixSlice[G]) Pop() (e *G, err error) {
	if m.IsEmpty() {
		return nil, io.EOF
	}
	e = m.array[m.head]
	m.array[m.head] = nil // 帮助GC回收内存
	m.head = (m.head + 1) % m.cap
	m.size--
	return
}

// 辅助函数：返回两数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 辅助函数：生成索引序列
func generateIndices(start, end int, asc bool) []int {
	if asc {
		indices := make([]int, end-start)
		for i := range indices {
			indices[i] = start + i
		}
		return indices
	} else {
		indices := make([]int, end-start)
		for i := range indices {
			indices[i] = end - 1 - i
		}
		return indices
	}
}

// 按队列顺序获取元素（先进先出）
func (m *FixSlice[G]) Queue(begin, limit int, orderAsc ...bool) []*G {
	return m.getElements(begin, limit, orderAsc, func(idx int) uint32 {
		return (m.head + uint32(idx)) % m.cap
	})
}

// 按栈顺序获取元素（后进先出）
func (m *FixSlice[G]) Stack(begin, limit int, orderAsc ...bool) []*G {
	return m.getElements(begin, limit, orderAsc, func(idx int) uint32 {
		return (m.tail - uint32(idx) - 1 + m.cap) % m.cap
	})
}

// 公共实现函数
func (m *FixSlice[G]) getElements(begin, limit int, orderAsc []bool, posFunc func(int) uint32) []*G {
	count := m.Len()
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

	asc := true
	if len(orderAsc) > 0 {
		asc = orderAsc[0]
	}

	indices := generateIndices(begin, end, asc)
	rows := make([]*G, len(indices))

	for i, idx := range indices {
		pos := posFunc(idx)
		rows[i] = m.array[pos]
	}

	return rows
}