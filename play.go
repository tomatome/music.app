// play.go
package main

import (
	"errors"
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

type ACTIVE_STATUS int

const (
	ACTIVE_STOPPED = iota
	ACTIVE_PLAYING
	ACTIVE_STALLED
	ACTIVE_PAUSED
)

type Music struct {
	curHandle uintptr
	isPlay    bool
	curUrl    string
	err       error
}

func initBass() error {
	bInit := bass.NewProc("BASS_Init")
	ret, _, _ := bInit.Call(uintptr(device), uintptr(freq), 0, 0, 0)

	if int(ret) != 0 {
		fmt.Println("initBass: ok...")
	} else {
		fmt.Println("initBass: fail...")
	}

	return nil
}

func (m *Music) checkError() {
	code, _, _ := BError.Call()
	ret := int(code)
	switch ret {
	case 0:
	case 2:
		m.err = errors.New("文件链接已失效")
		MsgBox("文件链接已失效")
	default:
		fmt.Println("checkError ret:", ret)
	}
}

func (m *Music) play(url string) error {
	if m.isActive() == ACTIVE_PLAYING {
		m.stop()
	}

	var (
		stream *syscall.LazyProc
		handle uintptr
		err    error
	)

	if strings.HasPrefix(url, "http") {
		stream = bass.NewProc("BASS_StreamCreateURL")
		handle, _, err = stream.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(url))),
			0,
			0x80000000,
			0,
			0)
		if err != nil {
			fmt.Println("streamURL:", err)
			//return err
		}
		m.checkError()
	} else {
		stream = bass.NewProc("BASS_StreamCreateFile")
		handle, _, err = stream.Call(0,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(url))),
			0,
			0,
			0x80000000)
		if err != nil {
			fmt.Println("streamFile:", err)
			//return err
		}
	}

	play := bass.NewProc("BASS_ChannelPlay")
	_, _, err = play.Call(handle, uintptr(1))
	if err != nil {
		fmt.Println("Play:", err)
		//return err
	}
	m.checkError()

	m.curHandle = handle
	m.curUrl = url
	m.isPlay = true

	return m.err
}

func (m *Music) stop() {
	stop := bass.NewProc("BASS_ChannelStop")
	_, _, err := stop.Call(m.curHandle)
	if err != nil {
		fmt.Println("Stop:", err)
	}
}

func (m *Music) isActive() ACTIVE_STATUS {
	if m.curHandle == uintptr(0) {
		return ACTIVE_STOPPED
	}
	isActive := bass.NewProc("BASS_ChannelIsActive")
	ret, _, err := isActive.Call(m.curHandle)
	if err != nil {
		fmt.Println("isActive:", err)
	}
	status := ACTIVE_STATUS(ret)
	fmt.Println("status:", status)
	return status
}
