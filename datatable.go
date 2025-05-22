/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-22 13:25:46
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 16:46:03
 * @FilePath: \gui\datatable.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gioui.org/text"
)

var _ Tabler = (*DataTable[int])(nil)

type Finder interface {
	Find(text string, colIdx int) (err error)
	ResetFind()
}

type Column struct {
	Width     int    //宽度
	Title     string // 列名
	Key       string //源列名
	Alignment text.Alignment
	Ids       []int
	cb        func(colIdx int) //单击标题时回调函数
}

func NewColumn(title string, width int, alig ...text.Alignment) *Column {
	alignment := text.Start
	if len(alig) > 0 {
		alignment = alig[0]
	}
	return &Column{Width: width, Title: title, Alignment: alignment}
}

func (m *Column) SetOnClick(cb func(colIdx int)) {
	m.cb = cb
}

type ItemText func(record any, row, col int) string

type DataTable[G any] struct {
	list     []*G
	root     []*G //原始数据不会变化
	cols     []*Column
	textFunc ItemText
}

// list 数组
// cols 可为空
// textFunc 回调函数
func NewDataTable[G any](list []*G, cols []*Column, textFunc ItemText) *DataTable[G] {
	m := new(DataTable[G])
	m.list = list
	m.root = list
	m.Init(cols, textFunc)
	return m
}

func (m *DataTable[G]) UpdateList(list []*G) {
	m.root = list
	m.list = list
}

func (m *DataTable[G]) GetList() []*G {
	return m.root
}

func (m *DataTable[G]) Init(cols []*Column, textFunc ItemText) {
	if cols == nil && textFunc == nil {
		cols, textFunc = MakeColumns(m.root, nil)
	}

	if cols != nil {
		m.cols = cols
	}
	if textFunc != nil {
		m.textFunc = textFunc
	}
}

// 获取标题
func (m *DataTable[G]) GetTitle(i int) (title string) {
	return m.cols[i].Title
}

// 获取列宽度
func (m *DataTable[G]) GetColumnWitdh(i int) (width float32) {
	return float32(m.cols[i].Width)
}

func (m *DataTable[G]) GetRow(row int) any {
	list := m.list
	n := len(list)
	if row >= n {
		return nil
	}

	return list[row]
}

// 获取列属性
func (m *DataTable[G]) GetColumn(i int) *Column {
	return m.cols[i]
}

// 获取列个数
func (m *DataTable[G]) GetColumnCount() (count int) {
	return len(m.cols)
}

// 获取行数
func (m *DataTable[G]) Size() (size int) {
	return len(m.list)
}

// 获取单元格数据
func (m *DataTable[G]) GetItemText(r any, row int, col int) (text string) {
	record := m.list[row]
	return m.textFunc(record, row, col)
}

func (m *DataTable[G]) ResetFind() {
	m.list = m.root
}

func (m *DataTable[G]) Find(text string, colIdx int) (err error) {
	text = strings.TrimSpace(text)
	if text == "" {
		err = errors.New("查找的字符串不能是空")
		return
	}

	list := reflect.ValueOf(m.root)
	list = reflect.Indirect(list)
	if !(list.Kind() == reflect.Slice || list.Kind() == reflect.Array) {
		err = errors.New("内部错误, 关联数组错误")
		return
	}

	size := m.Size()
	if size == 0 {
		return
	}

	idx := 0

	newList := reflect.New(list.Type()).Elem()
	newList.Set(reflect.MakeSlice(list.Type(), list.Len(), list.Cap()))

	defer func() {
		if idx >= 0 {
			newList.SetLen(idx)
			m.list = newList.Interface().([]*G)
		} else {
			err = errors.New("not found")
		}
	}()

	for i := 0; i < size; i++ {
		row := m.GetRow(i)
		if strings.Contains(m.GetItemText(row, i, colIdx), text) {
			newList.Index(idx).Set(list.Index(i))
			idx++
			continue
		}
	}
	return
}

func (m *DataTable[G]) GetRowString(record any, row int, buf io.Writer) {
	n := m.GetColumnCount()
	for i := 0; i < n; i++ {
		if i != 0 {
			buf.Write([]byte(","))
		}
		buf.Write([]byte(m.textFunc(record, row, i)))
	}
}

func MakeColumns(list any, mapColumn map[string]*Column) (cols []*Column, f ItemText) {
	key := map[string][]int{}
	t := reflect.TypeOf(list)
	t = t.Elem()
	var names []string
	var shortNames []string
	MapNameID("", key, nil, t, true, &names, nil, &shortNames)

	cols = make([]*Column, len(key)+1)
	c := &Column{
		Title: "序号",
		Width: 75,
	}

	cols[0] = c

	for j, name := range names {
		if mapColumn != nil {
			column, ok := mapColumn[name]
			if ok && column != nil {
				cols[j+1] = column
				cols[j+1].Ids = key[name]
				continue
			}
		}
		col := Column{Title: shortNames[j]}
		col.Ids = key[name]
		col.SetOnClick(func(colIdx int) {
			//m.win.Log(col.Title)
		})
		cols[j+1] = &col
	}
	f = func(row any, rowIdx, col int) (text string) {
		if col == 0 {
			return strconv.Itoa(rowIdx + 1)
		}
		v := reflect.ValueOf(row)
		v = reflect.Indirect(v)

		v = v.FieldByIndex(cols[col].Ids)
		if strings.HasSuffix(cols[col].Key, "time") {
			if v.CanInt() {
				return time.Unix(v.Int(), 0).Format("2006-01-02 15:04:05")
			}
			return fmt.Sprint(v)
		} else if v.CanFloat() {
			return strconv.FormatFloat(v.Float(), 'f', 3, 64)
		}

		return fmt.Sprint(v)
	}
	return
}

/**
 * @description: struct 字段名转为ids
 * @param {string} prefix
 * @param {map[string][]int} key
 * @param {[]int} paths
 * @param {reflect.Type} t
 * @param {bool} toLower 是否转为小写
 * @return {*}
 */
func MapNameID(prefix string, mapKeyId map[string][]int, paths []int, t reflect.Type, toLower bool, keys ...*[]string) {
	var names *[]string
	var srcNames *[]string
	var shortName *[]string
	if len(keys) > 0 {
		names = keys[0]
	}
	if len(keys) > 1 {
		srcNames = keys[1]
	}
	if len(keys) > 2 {
		shortName = keys[2]
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	//log.Println(prefix + t.String())
	for i := 0; i < t.NumField(); i++ {
		newPaths := make([]int, len(paths))
		copy(newPaths, paths)
		newPaths = append(newPaths, i)
		fieldType := t.Field(i)

		typ := fieldType.Type

		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		//log.Println(prefix+fieldType.Name, typ.Kind())

		if fieldType.Anonymous {
			if typ.Kind() == reflect.Struct {
				MapNameID(prefix, mapKeyId, newPaths, typ, toLower, names, srcNames, shortName)
			} else {
				name := prefix + fieldType.Name
				if srcNames != nil {
					*srcNames = append(*srcNames, name)
				}
				if shortName != nil {
					*shortName = append(*shortName, fieldType.Name)
				}
				if toLower {
					name = strings.ToLower(name)
				}
				mapKeyId[name] = newPaths
				if names != nil {
					*names = append(*names, name)
				}

			}
			continue
		}
		if typ.Kind() == reflect.Struct {
			tag := fieldType.Tag.Get("json")
			if tag == "-" {
				continue
			}
			if fieldType.IsExported() {
				MapNameID(prefix+fieldType.Name+".", mapKeyId, newPaths, typ, toLower, names)
			}
			continue
		}

		if fieldType.IsExported() {
			name := prefix + fieldType.Name
			if srcNames != nil {
				*srcNames = append(*srcNames, name)
			}
			if shortName != nil {
				*shortName = append(*shortName, fieldType.Name)
			}

			if toLower {
				name = strings.ToLower(name)
			}
			mapKeyId[name] = newPaths
			if names != nil {
				*names = append(*names, name)
			}
		}
	}
}
