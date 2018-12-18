// mui.go
package main

import (
	"encoding/json"
	"fmt"
	"music/madoka"
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func initMenuItem() []MenuItem {
	var openAction, showAboutBoxAction *walk.Action
	var recentMenu *walk.Menu
	return []MenuItem{
		Menu{
			Text: "&文件",
			Items: []MenuItem{
				Action{
					AssignTo:    &openAction,
					Text:        "&打开",
					Enabled:     Bind("enabledCB.Checked"),
					Visible:     Bind("!openHiddenCB.Checked"),
					Shortcut:    Shortcut{walk.ModControl, walk.KeyO},
					OnTriggered: mw.openfiles,
				},
				Menu{
					AssignTo: &recentMenu,
					Text:     "Recent",
				},
				Separator{},
				Action{
					Text:        "E&xit",
					OnTriggered: func() { mw.Close() },
				},
			},
		},
		Menu{
			Text: "&排行榜",
			Items: []MenuItem{
				Action{
					Text: "云音乐热歌榜",
					OnTriggered: func() {
						mw.setTitle("云音乐热歌榜")
						mw.model.getPlayListById("3778678")
						mw.model.Publish()
					},
				},
				Action{
					Text: "云音乐新歌榜",
					OnTriggered: func() {
						mw.setTitle("云音乐新歌榜")
						mw.model.getPlayListById("3779629")
						mw.model.Publish()
					},
				},
				Action{
					Text: "华语金曲榜",
					OnTriggered: func() {
						mw.setTitle("华语金曲榜")
						mw.model.getPlayListById("4395559")
						mw.model.Publish()
					},
				},
				Action{
					Text: "中国TOP排行榜（内地榜）",
					OnTriggered: func() {
						mw.setTitle("中国TOP排行榜（内地榜）")
						mw.model.getPlayListById("64016")
						mw.model.Publish()
					},
				},
				Action{
					Text: "中国TOP排行榜（港台榜）",
					OnTriggered: func() {
						mw.setTitle("中国TOP排行榜（港台榜）")
						mw.model.getPlayListById("112504")
						mw.model.Publish()
					},
				},
				Action{
					Text: "云音乐飙升榜",
					OnTriggered: func() {
						mw.setTitle("云音乐飙升榜")
						mw.model.getPlayListById("19723756")
						mw.model.Publish()
					},
				},
				Action{
					Text: "网易原创歌曲榜",
					OnTriggered: func() {
						mw.setTitle("网易原创歌曲榜")
						mw.model.getPlayListById("2884035")
						mw.model.Publish()
					},
				},
			},
		},
		Menu{
			Text:  "&歌单",
			Items: mw.initMusicSheet(),
		},
		Menu{
			Text:  "&电台",
			Items: mw.initFmList(),
		},
		Menu{
			Text: "&帮助",
			Items: []MenuItem{
				Action{
					AssignTo: &showAboutBoxAction,
					Text:     "关于",
					OnTriggered: func() {
					},
				},
			},
		},
	}
}

func (mw *MyMainWindow) initMusicSheet() []MenuItem {
	item := make([]MenuItem, 0, 30)
	resp, _ := madoka.PlayList("全部", "hot", 1, 29)
	//fmt.Println(resp)
	/*a := Action{Text: "歌单440103454",
		OnTriggered: func() {
			mw.setTitle("经典怀旧")
			mw.model.getPlayListById("440103454")
			mw.model.Publish()
		},
	}
	item = append(item, a)*/

	var f map[string]interface{}
	json.Unmarshal([]byte(resp), &f)
	list := f["playlists"].([]interface{})
	for _, lv := range list {
		v := lv.(map[string]interface{})
		t := int(v["id"].(float64))
		id := strconv.Itoa(t)
		a := Action{Text: fmt.Sprintf("%6s", v["name"].(string)),
			OnTriggered: func() {
				mw.setTitle(v["name"].(string))
				mw.model.getPlayListById(id)
				mw.model.Publish()
			},
		}
		item = append(item, a)
	}
	return item
}

func (mw *MyMainWindow) initFmList() []MenuItem {
	item := make([]MenuItem, 0, 300)
	resp, _ := madoka.FmCatalogue()
	//fmt.Println(resp)
	var f map[string]interface{}
	json.Unmarshal([]byte(resp), &f)
	list := f["categories"].([]interface{})
	for _, lv := range list {
		v := lv.(map[string]interface{})
		t := int(v["id"].(float64))
		id := strconv.Itoa(t)
		a := Action{Text: v["name"].(string),
			OnTriggered: func() {
				mw.setTitle(v["name"].(string))
				mw.model.getFmInfoByName(v["name"].(string), id)
				mw.model.Publish()
			},
		}
		item = append(item, a)
	}
	return item
}

type Species struct {
	Id   int
	Name string
}

func KnownSpecies() []*Species {
	return []*Species{
		{1, "单曲"},
		{100, "歌手"},
		{1000, "歌单"},
	}
}

func setSearchType() []string {
	return []string{
		"单曲",
		"歌手",
		"歌单",
	}
}

func getSearchTypeByName(name string) string {
	switch name {
	case "歌手":
		return "100"
	case "歌单":
		return "1000"
	default:
	}
	return "1"
}

func initSearchBox() GroupBox {
	return GroupBox{
		Layout: HBox{},
		Children: []Widget{
			LineEdit{
				AssignTo: &mw.searchBox,
				OnKeyDown: func(key walk.Key) {
					//fmt.Println(key, walk.KeyReturn)
					if key == walk.KeyReturn {
						mw.clicked()
					}
				},
			},
			ComboBox{
				AssignTo:     &mw.comBox,
				CurrentIndex: 0,
				Model:        setSearchType(),
			},
			PushButton{
				Text:      "搜索",
				OnClicked: mw.clicked,
			},
		},
	}
}

func initTableView() TableView {
	return TableView{
		AssignTo: &mw.tv,
		//AlternatingRowBGColor: walk.RGB(239, 239, 239),
		CheckBoxes:       true,
		ColumnsOrderable: true,
		MultiSelection:   true,
		Columns: []TableViewColumn{
			{Title: "序号"},
			{Title: "歌名", Width: 100, Alignment: AlignCenter},
			{Title: "作者", Alignment: AlignCenter},
			{Title: "链接", Alignment: AlignCenter, Width: 400},
		},
		StyleCell: func(style *walk.CellStyle) {
			item := mw.model.items[style.Row()]

			if item.Checked {
				if style.Row()%2 == 0 {
					style.BackgroundColor = walk.RGB(159, 215, 255)
				} else {
					style.BackgroundColor = walk.RGB(143, 199, 239)
				}
			}

			switch style.Col() {
			case 1:
				/*if canvas := style.Canvas(); canvas != nil {
					bounds := style.Bounds()
					bounds.X += 2
					bounds.Y += 2
					bounds.Width = int((float64(bounds.Width) - 4) / 5 * float64(len(item.Name)))
					bounds.Height -= 4
					canvas.DrawBitmapPartWithOpacity(barBitmap, bounds, walk.Rectangle{0, 0, 100 / 5 * len(item.Name), 1}, 127)

					bounds.X += 4
					bounds.Y += 2
					canvas.DrawText(item.Name, mw.tv.Font(), 0, bounds, walk.TextLeft)
				}*/

			case 2:
				/*if item.Baz >= 900.0 {
					style.TextColor = walk.RGB(0, 191, 0)
					//style.Image = goodIcon
				} else if item.Baz < 100.0 {
					style.TextColor = walk.RGB(255, 0, 0)
					//style.Image = badIcon
				}*/

			case 3:
				/*if item.Quux.After(time.Now().Add(-365 * 24 * time.Hour)) {
					style.Font = boldFont
				}*/
			}
		},
		Model: mw.model,
		OnItemActivated: func() {
			i := mw.tv.CurrentIndex()
			if i < 0 {
				return
			}
			fmt.Printf("OnItemActivated: index=%d, Duration=%d\n", i, mw.model.items[i].Duration)
			Url := mw.model.items[i].Url
			mw.play(i, Url)
		},
	}
}

func initControl() GroupBox {
	return GroupBox{
		Layout: HBox{},
		Children: []Widget{
			PushButton{
				AssignTo: &mw.ctrlPanel,
				//ImageAboveText: true,
				Image:     "image/play.png",
				Text:      " ",
				OnClicked: mw.musicClicked,
			},
			PushButton{
				AssignTo:  &mw.pModel,
				Text:      "单曲",
				OnClicked: mw.setPlayModel,
			},
			GroupBox{
				Layout:        VBox{},
				StretchFactor: 2,
				Children: []Widget{
					GroupBox{
						Layout: HBox{},
						Children: []Widget{
							Label{
								AssignTo: &mw.lname,
								Text:     "未知",
							},
							Label{
								AssignTo:           &mw.ltime,
								Text:               "00:00",
								RightToLeftReading: true,
							},
						},
					},
					ProgressBar{
						AssignTo: &mw.pb,
						Value:    0,
					},
				},
			},
			/*PushButton{
				Background:     SolidColorBrush{Color: walk.RGB(255, 255, 255)},
				Image:          "music.ico",
				ImageAboveText: true,
				OnClicked:      mw.showList,
				StretchFactor:  1,
			},*/
		},
	}
}
