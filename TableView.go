// TableView.go
package main

import (
	"sort"

	"github.com/lxn/walk"
)

type Foo struct {
	Index    int
	Id       string
	Name     string
	Artist   string
	Url      string
	Duration int
	Checked  bool
}

type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo
}

func NewFooModel() *FooModel {
	m := new(FooModel)
	m.ResetRows()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *FooModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Name

	case 2:
		return item.Artist

	case 3:
		return item.Url
	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *FooModel) Checked(row int) bool {
	return m.items[row].Checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FooModel) SetChecked(row int, checked bool) error {
	m.items[row].Checked = checked

	return nil
}

// Called by the TableView to sort the model.
func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.Name < b.Name)

		case 2:
			return c(a.Artist < b.Artist)

		case 3:
			return c(a.Artist < b.Artist)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

type mModel struct {
	Playlists []*Playlists
}
type Playlists struct {
	Name string
	Id   int
}

func (m *FooModel) ResetRows() {
	m.items = make([]*Foo, 0, 100)
	/*a, _ := madoka.PlayList("全部", "hot", 0, 2)
	//fmt.Println("a:", a)
	var md mModel
	json.Unmarshal([]byte(a), &md)
	fmt.Println(md.Playlists[1].Name)
	f := new(Foo)
	f.Name = md.Playlists[1].Name
	f.Id = strconv.Itoa(md.Playlists[1].Id)
	//fmt.Println("f.Id:", f.Id)
	m.items = append(m.items, f)*/

	for i := range m.items {
		m.items[i].Index = i
	}

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}
