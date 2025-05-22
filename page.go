package gui

type Page struct {
	active  bool
	title   string
	id      string //必须唯一
	content Contenter
	w       *Window
}

func (m *Page) Setup(w *Window, id string, title string) {
	m.w = w
	m.id = id
	m.title = title
}

func (m *Page) ID() string {
	return m.id
}
func (m *Page) Title() string {
	return m.title
}

func (m *Page) Actived() bool {
	return m.active
}

func (m *Page) Show() *Page {
	m.active = true
	return m
}

func (m *Page) Hide() *Page {
	m.active = false
	return m
}

func (m *Page) Parent() *Window {
	return m.w
}

func (m *Page) SetContent(content Contenter) *Page {
	m.content = content
	return m
}

func (m *Page) Layout(gtx C) D {
	if m.content == nil {
		return D{}
	}
	return m.content.Layout(gtx)
}
