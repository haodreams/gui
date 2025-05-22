package gui

import "errors"

//管理页面的路由
type Route struct {
	history []*Page
	w       *Window
	pages   []*Page
	mapPage map[string]*Page
	curPage *Page
}

func NewRoute(w *Window) *Route {
	m := new(Route)
	m.mapPage = make(map[string]*Page)
	m.w = w
	return m
}

func (m *Route) To(to *Page) {
	if m.curPage == nil {
		m.curPage = to
		m.curPage.Show()
		return
	}
	if m.curPage.ID() == to.ID() {
		return
	}
	m.history = append(m.history, m.curPage)
	m.curPage = to
	m.curPage.Hide()
	to.Show()
}

func (m *Route) Back() {
	if len(m.history) == 0 {
		return
	}
	to := m.history[len(m.history)-1]
	m.history = m.history[:len(m.history)-1]
	m.To(to)
}

func (m *Route) AddPage(page *Page) (err error) {
	if m.curPage == nil {
		m.curPage = page
		m.curPage.Show()
	}
	p, ok := m.mapPage[page.ID()]
	if ok && p != nil {
		return errors.New("page " + p.ID() + " is existed")
	}
	m.pages = append(m.pages, page)
	return
}

func (m *Route) Layout(gtx C) D {
	if m.curPage == nil {
		return D{}
	}
	return m.curPage.Layout(gtx)
}
