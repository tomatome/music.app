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

const (
	Title = "音乐播放器"
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

func init() {
	mw = &MyMainWindow{
		pbchan: make(chan int, 1),
	}

	mw.music = new(Music)
	mw.model = NewFooModel()
}

func main() {
	initBass()

	var sbi *walk.StatusBarItem
	//boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)
	barBitmap, err := walk.NewBitmap(walk.Size{100, 1})
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	MW := MainWindow{
		AssignTo:  &mw.MainWindow,
		Title:     Title,
		Icon:      "music.ico",
		MinSize:   Size{800, 600},
		MenuItems: initMenuItem(),
		Layout:    VBox{},
		Children: []Widget{
			initSearchBox(),
			initTableView(),
			initControl(),
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
	menu       *walk.Menu
	tv         *walk.TableView
	model      *FooModel
	comBox     *walk.ComboBox
	searchBox  *walk.LineEdit
	music      *Music
	ctrlPanel  *walk.PushButton
	pModel     *walk.PushButton
	lname      *walk.Label
	ltime      *walk.Label
	pb         *walk.ProgressBar
	pbchan     chan int
	musicModel int
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

func (mw *MyMainWindow) setTitle(s string) {
	mw.SetTitle(Title + " ● " + s)
}
func (mw *MyMainWindow) setPlayModel() {
	if !mw.music.isPlay {
		return
	}
	if mw.pModel.Text() == "单曲" {
		mw.pModel.SetText("循环")
		mw.musicModel = 1
	} else {
		mw.pModel.SetText("单曲")
		mw.musicModel = 0
	}

	/*if mw.musicModel == 1 {
		for i := mw.tv.CurrentIndex(); i < len(mw.model.items); {
			i++
			if i > len(mw.model.items) {
				i = 0
			}
			Url := mw.model.items[i].Url
			mw.play(i, Url)
			mw.tv.SetCurrentIndex(i)
			mw.model.items[i].Checked = true
			num := 400
			for {
				if mw.music.isActive() != ACTIVE_PLAYING {
					break
				}
				if num == 0 {
					break
				}
				num--
				time.Sleep(1 * time.Second)
			}
		}
	}*/

	//if mw.music.isActive() == ACTIVE_PAUSED
}

func (mw *MyMainWindow) play(i int, url string) {
	mw.music.curUrl = url
	mw.lname.SetText(mw.model.items[i].Name + " --- " + mw.model.items[i].Artist)
	mw.model.items[i].Checked = true
	mw.pb.SetRange(0, mw.model.items[i].Duration)
	mw.MusicStart()
}

func (mw *MyMainWindow) musicClicked() {
	fmt.Println("musicClicked:", mw.music.isPlay)
	if mw.music.curUrl == "" {
		return
	}
	if mw.music.isPlay {
		mw.MusicStop()
	} else {
		mw.MusicStart()
	}
}
func (mw *MyMainWindow) clicked() {
	req := mw.searchBox.Text()
	if strings.TrimSpace(req) == "" {
		mw.model.ResetRows()
		return
	}

	t := getSearchTypeByName(mw.comBox.Text())
	resp, _ := madoka.Search(req, t, 1, 20)
	fmt.Println(resp)
	var ms MSearch
	json.Unmarshal([]byte(resp), &ms)
	fmt.Println("all:", len(ms.Result.Songs))
	mw.model.items = mw.model.items[0:0]

	ids := make([]string, 0, len(ms.Result.Songs))
	allMap := make(map[int]SongList)
	for _, s := range ms.Result.Songs {
		ids = append(ids, strconv.Itoa(s.Id))
		allMap[s.Id] = s
	}

	a, _ := madoka.Download("["+strings.Join(ids, ",")+"]", "320000")
	var m SongInfo
	json.Unmarshal([]byte(a), &m)
	//fmt.Println(a)
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

	mw.model.Publish()
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
	walk.MsgBox(mw, Title, message, walk.MsgBoxIconWarning)
}

const (
	MUSIC_STOP = iota
	MUSIC_START
	MUSIC_RESTART
)

func (mw *MyMainWindow) MusicStart() {
	mw.music.play(mw.music.curUrl)
	play, _ := walk.NewImageFromFile("image/play.png")
	mw.ctrlPanel.SetImage(play)
	mw.pbchan <- mw.pb.MaxValue()
}

func (mw *MyMainWindow) MusicStop() {
	mw.music.stop()
	pause, _ := walk.NewImageFromFile("image/pause.png")
	mw.ctrlPanel.SetImage(pause)
	mw.pbchan <- MUSIC_STOP
}
func (mw *MyMainWindow) MusicRestart() {
	mw.MusicStop()
	mw.MusicStart()
}
func (mw *MyMainWindow) startProgressBar() {
	num := -1
	for {
		select {

		case start := <-mw.pbchan:
			if start == MUSIC_STOP {
				mw.pb.SetValue(0)
				mw.ltime.SetText("0 / " + strconv.Itoa(mw.pb.MaxValue()))
				num = -1
			} else {
				mw.pb.SetValue(0)
				num = 0
			}
			//fmt.Println(start)
		case <-time.After(1 * time.Second):
		}

		if num != -1 {
			if num < mw.pb.MaxValue() {
				num++
				mw.pb.SetValue(num)
				mw.ltime.SetText(strconv.Itoa(num) + " / " + strconv.Itoa(mw.pb.MaxValue()))
			}
		}
	}
}
func (mw *MyMainWindow) showList() {

}
