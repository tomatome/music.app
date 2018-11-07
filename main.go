package main

import (
	"encoding/json"
	"fmt"
	"log"
	"music/madoka"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Detail struct {
	Id   int
	Url  string
	Code int
	Type string
}
type SongInfo struct {
	Data []Detail
}

var mw *MyMainWindow

func main() {
	initBass()
	mw = &MyMainWindow{pbchan: make(chan int, 1)}
	mw.music = new(Music)
	var openAction, showAboutBoxAction *walk.Action
	var recentMenu *walk.Menu
	var sbi *walk.StatusBarItem
	//boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)
	barBitmap, err := walk.NewBitmap(walk.Size{100, 1})
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()
	mw.model = NewFooModel()
	MW := MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "音乐播放器",
		Icon:     "music.ico",
		MinSize:  Size{800, 600},
		MenuItems: []MenuItem{
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
				Text: "&Help",
				Items: []MenuItem{
					Action{
						AssignTo: &showAboutBoxAction,
						Text:     "预览",
						OnTriggered: func() {

						},
					},
				},
			},
		},
		Layout: VBox{},

		Children: []Widget{
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.searchBox,
					},
					PushButton{
						Text:      "搜索",
						OnClicked: mw.clicked,
					},
				},
			},
			TableView{
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
					fmt.Printf("OnItemActivated: %v\n", i)
					Url := mw.model.items[i].Url
					mw.music.play(Url)
					mw.lname.SetText(mw.model.items[i].Name + " --- " + mw.model.items[i].Artist)
					mw.pbchan <- mw.model.items[i].Duration / 1000
					mw.model.items[i].Checked = true
				},
			},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &mw.ctrlPanel,
						//Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						Image:          "image/play.png",
						ImageAboveText: true,
						//MaxSize: Size{72, 72},
						//MinSize: Size{72, 72},
						Text: "  ",
						//StretchFactor:  2,
						OnClicked: mw.musicClicked,
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
					PushButton{
						Background:     SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						Image:          "music.ico",
						ImageAboveText: true,
						OnClicked:      mw.showList,
						StretchFactor:  1,
					},
				},
			},
		},
		StatusBarItems: []StatusBarItem{
			StatusBarItem{
				AssignTo: &sbi,
				Text:     "click",
				Width:    80,
			},
		},
	}
	go mw.startProgressBar()
	if _, err := MW.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("ok")
}

type MyMainWindow struct {
	*walk.MainWindow
	menu      *walk.Menu
	tv        *walk.TableView
	model     *FooModel
	searchBox *walk.LineEdit
	music     *Music
	ctrlPanel *walk.PushButton
	lname     *walk.Label
	ltime     *walk.Label
	pb        *walk.ProgressBar
	pbchan    chan int
}
type MSearch struct {
	Result Song
}
type Song struct {
	Songs []SongList
}
type SongList struct {
	Id       int
	Name     string
	Artists  []Artist
	Duration int
}
type Artist struct {
	Name string
}

func (mw *MyMainWindow) musicClicked() {
	fmt.Println(mw.music.isPlay)
	if mw.music.file == "" {
		return
	}
	if mw.music.isPlay {
		mw.music.stop()
		pause, _ := walk.NewImageFromFile("image/pause.png")
		mw.ctrlPanel.SetImage(pause)
		mw.pbchan <- -1
	} else {
		mw.music.play(mw.music.file)
		play, _ := walk.NewImageFromFile("image/play.png")
		mw.ctrlPanel.SetImage(play)
		mw.pbchan <- 0
	}
}
func (mw *MyMainWindow) clicked() {
	req := mw.searchBox.Text()
	resp, _ := madoka.Search(req, "1", 1, 100)
	fmt.Println(resp)
	var ms MSearch
	json.Unmarshal([]byte(resp), &ms)
	fmt.Println("all:", len(ms.Result.Songs))
	mw.model.ResetRows()
	ids := make([]string, 0, len(ms.Result.Songs))
	allMap := make(map[int]SongList)
	for _, s := range ms.Result.Songs {
		ids = append(ids, strconv.Itoa(s.Id))
		allMap[s.Id] = s
	}
	fmt.Println(strings.Join(ids, ","))
	a, _ := madoka.Download("["+strings.Join(ids, ",")+"]", "320000")
	var m SongInfo
	json.Unmarshal([]byte(a), &m)
	fmt.Println(a)
	for i, data := range m.Data {
		if data.Code != 200 {
			continue
		}
		s := allMap[data.Id]
		f := new(Foo)
		f.Name = s.Name
		f.Id = strconv.Itoa(s.Id)
		f.Url = data.Url
		f.Artist = s.Artists[0].Name
		f.Duration = s.Duration / 1000
		f.Index = i
		//fmt.Println("f.Id:", f.Id)
		mw.model.items = append(mw.model.items, f)
	}

	mw.model.PublishRowsReset()
}
func (mw *MyMainWindow) openfiles() {
	dlg := new(walk.FileDialog)
	dlg.FilePath = ""
	dlg.Title = "Select File"
	dlg.Filter = "Exe files (*.mp3)|*.mp3|All files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		//mw.preview.SetText("Error : File Open\r\n")
		return
	} else if !ok {
		//mw.preview.SetText("Cancel\r\n")
		return
	}
	path := dlg.FilePath
	s := fmt.Sprintf("Select : %s\r\n", path)
	fmt.Println(s)
	mw.music.play(path)

}
func search(text, word string) (result []int) {
	result = []int{}
	i := 0
	for j, _ := range text {
		if strings.HasPrefix(text[j:], word) {
			log.Print(i)
			result = append(result, i)
		}
		i += 1
	}
	return
}

func MsgBox(message string) {
	walk.MsgBox(mw, "音乐播放器", message, walk.MsgBoxIconWarning)
}

func (mw *MyMainWindow) startProgressBar() {
	num := -1
	for {
		select {

		case start := <-mw.pbchan:
			fmt.Println(start)
			if start == -1 {
				num = -1
				mw.pb.SetValue(0)
				mw.ltime.SetText("00:00:00/00:00:" + strconv.Itoa(mw.pb.MaxValue()))
			} else if start == 0 {
				mw.pb.SetValue(0)
				num = 0
			} else {
				mw.pb.SetValue(0)
				mw.pb.SetRange(0, start)
				num = 0
			}
		case <-time.After(1 * time.Second):
			//fmt.Println("超时了")
			//default:
			//fmt.Println("default")
		}

		if num != -1 {
			if num < mw.pb.MaxValue() {
				num++
				mw.pb.SetValue(num)
				mw.ltime.SetText("00:00:" + strconv.Itoa(num) + "/00:00:" + strconv.Itoa(mw.pb.MaxValue()))
			}
		}
	}
}
func (mw *MyMainWindow) showList() {

}
