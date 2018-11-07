// play.go
package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	bass   = syscall.NewLazyDLL("bass24/x64/bass.dll")
	BError = bass.NewProc("BASS_ErrorGetCode")
	device = -1
	freq   = 44100
)

type Music struct {
	curHandle uintptr
	isPlay    bool
	file      string
}

func initBass() error {
	bInit := bass.NewProc("BASS_Init")
	ret, _, _ := bInit.Call(uintptr(device), uintptr(freq), 0, 0, 0)

	if int(ret) != 0 {
		fmt.Println("init: ok...")
	} else {
		fmt.Println("init: fail...")
	}

	return nil
}

func checkError() {
	code, _, _ := BError.Call()
	ret := int(code)
	switch ret {
	case 0:
		fmt.Println("ok")
	case 2:
		MsgBox("文件链接已失效")
	default:
		fmt.Println("checkError ret:", ret)
	}
}

func (m *Music) play(file string) error {
	if m.curHandle != uintptr(0) {
		m.stop()
	}
	var stream *syscall.LazyProc
	var source uintptr
	var err error
	if strings.HasPrefix(file, "http") {
		stream = bass.NewProc("BASS_StreamCreateURL")
		source, _, err = stream.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(file))),
			0,
			0x80000000,
			0,
			0)
		if err != nil {
			fmt.Println("streamURL:", err)
			//return err
		}
		checkError()
	} else {
		stream = bass.NewProc("BASS_StreamCreateFile")
		source, _, err = stream.Call(0,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(file))),
			0,
			0,
			0x80000000)
		if err != nil {
			fmt.Println("streamFile:", err)
			//return err
		}
	}

	play := bass.NewProc("BASS_ChannelPlay")
	_, _, err = play.Call(source, uintptr(1))
	if err != nil {
		fmt.Println("Play:", err)
		//return err
	}
	checkError()
	m.curHandle = source
	m.file = file
	m.isPlay = true

	return nil
}

func (m *Music) stop() {
	stop := bass.NewProc("BASS_ChannelStop")
	_, _, err := stop.Call(m.curHandle)
	if err != nil {
		fmt.Println("Stop:", err)
		//return
	}
	checkError()
	m.isPlay = false
}

func test() {
	//music_path := "https://m10.music.126.net/20181011142710/ab9b6aa96b069332b2af878cb5cb9b62/ymusic/07fa/a2a1/35ea/732937117d6d0a8c13a81bb40184662e.mp3"
	//music_path := "mu.mp3"
	//music_path := "http://www.170mv.com/kw/other.web.nn01.sycdn.kuwo.cn/resource/n3/44/77/2277939288.mp3"
	//initBass()
	//play(music_path)
}
