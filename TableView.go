// TableView.go
package main

import (
	"encoding/json"
	"fmt"
	"music/madoka"
	"sort"
	"strconv"
	"strings"

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
	m.items = make([]*Foo, 0, 100)
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

func (m *FooModel) getFmInfoByName(cat, id string) {
	m.items = m.items[0:0]
	a, _ := madoka.FmHotList(cat, id, 1, 30)
	fmt.Println(a)
	var f map[string]interface{}
	json.Unmarshal([]byte(a), &f)
	list := f["djRadios"].([]interface{})
	allMap := make(map[int]map[string]interface{})
	ids := make([]string, 0, len(list))
	for _, lv := range list {
		v := lv.(map[string]interface{})
		id := int(v["id"].(float64))
		ids = append(ids, strconv.Itoa(id))
		allMap[id] = v
	}
	a1, _ := madoka.Download("["+strings.Join(ids, ",")+"]", "320000")
	var ms SongInfo
	json.Unmarshal([]byte(a1), &ms)
	for i, data := range ms.Data {
		if data.Code != 200 {
			continue
		}
		v := allMap[data.Id]
		f := new(Foo)
		f.Name = v["name"].(string)
		f.Id = strconv.Itoa(int(v["id"].(float64)))
		f.Url = data.Url
		//f.Artist = v["artists"].([]interface{})[0].(map[string]interface{})["name"].(string)
		//f.Duration = int(v["duration"].(float64)) / 1000
		f.Index = i
		m.items = append(m.items, f)
	}
}

func (m *FooModel) getPlayListById(id string) {
	m.items = m.items[0:0]
	a, _ := madoka.PlayListDetail(id)
	var f map[string]interface{}
	json.Unmarshal([]byte(a), &f)
	list := f["result"].(map[string]interface{})["tracks"].([]interface{})
	//fmt.Println("list:", list)
	allMap := make(map[int]map[string]interface{})
	ids := make([]string, 0, len(list))
	for _, lv := range list {
		v := lv.(map[string]interface{})
		id := int(v["id"].(float64))
		ids = append(ids, strconv.Itoa(id))
		allMap[id] = v
	}
	a1, _ := madoka.Download("["+strings.Join(ids, ",")+"]", "320000")
	var ms SongInfo
	json.Unmarshal([]byte(a1), &ms)
	for i, data := range ms.Data {
		if data.Code != 200 {
			continue
		}

		v := allMap[data.Id]
		f := new(Foo)
		f.Name = v["name"].(string)
		f.Id = strconv.Itoa(int(v["id"].(float64)))
		f.Url = data.Url
		f.Artist = v["artists"].([]interface{})[0].(map[string]interface{})["name"].(string)
		f.Duration = int(v["duration"].(float64) / 1000)
		f.Index = i
		m.items = append(m.items, f)
	}
	fmt.Println(len(m.items))
}

func (m *FooModel) ResetRows() {
	m.getPlayListById("3778678")
	m.Publish()
}

func (m *FooModel) Publish() {
	m.PublishRowsReset()
	//m.Sort(m.sortColumn, m.sortOrder)
}
