// play.go
package main

import (
	"testing"
)

func TestPlayMusic(t *testing.T) {
	//music_path := "https://m10.music.126.net/20181011142710/ab9b6aa96b069332b2af878cb5cb9b62/ymusic/07fa/a2a1/35ea/732937117d6d0a8c13a81bb40184662e.mp3"
	music_path := "mu.mp3"
	//music_path := "http://www.170mv.com/kw/other.web.nn01.sycdn.kuwo.cn/resource/n3/44/77/2277939288.mp3"
	initBass()
	play(music_path)
}
