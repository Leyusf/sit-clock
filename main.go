package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var interval = 5
var a fyne.App
var restTime = interval
var flag = false
var music = false
var audioName = ""
var musicName = "默认"
var out = make(chan int)
var ch = make(chan int)
var audio = make([][]byte, 3)
var img = ToFile(clockImg)
var path string
var context *oto.Context
var lastIndex = -1

func readConfig(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	var buf [256]byte
	n, err := file.Read(buf[:])
	if err != nil {
		return ""
	}
	defer file.Close()
	strs := strings.Split(string(buf[:n]), "/")
	musicName = strs[len(strs)-1]
	return string(buf[:n])
}

func writeConfig(path, content string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write([]byte(content))
}

func main() {
	rand.Seed(time.Now().Unix())
	if runtime.GOOS == "linux" {
		path = "/home/config.cfg"
	} else if runtime.GOOS == "windows" {
		path = "C:\\Config\\clockConfig.cfg"
	}
	audioName = readConfig(path)
	audio[0] = ToFile(cai1)
	audio[1] = ToFile(cai2)
	audio[2] = ToFile(cai3)
	DemoClock()
}

func DemoClock() {
	a = app.New()

	CreateWindow()

	a.Run()
}

func ToFile(bs64String string) []byte {
	base64Bytes, err := base64.StdEncoding.DecodeString(bs64String)
	if err != nil {
		log.Fatal(err)
	}
	return base64Bytes
}

var filter = storage.NewExtensionFileFilter([]string{".mp3"})

func CreateWindow() {
	w := a.NewWindow("Clock")
	w.Resize(fyne.NewSize(500, 400))
	w.CenterOnScreen()

	newImg := fyne.NewStaticResource("clock", img)

	img := canvas.NewImageFromResource(newImg)
	img.FillMode = canvas.ImageFillOriginal
	content1 := container.New(layout.NewCenterLayout(), container.New(layout.NewGridWrapLayout(fyne.NewSize(80, 80)), img))

	clock := canvas.NewText("", color.Black)
	empty := canvas.NewText(musicName, color.Black)

	selectMusicBtn := widget.NewButton("选择音乐", func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			if read == nil {
				return
			}

			defer read.Close()
			audioName = read.URI().Path()
			empty.Text = read.URI().Name()
			empty.Refresh()
			writeConfig(path, audioName)
		}, w)
		openDialog.SetFilter(filter)
		openDialog.Show()
	})

	startBtn := widget.NewButton("开始", func() {
		Stop()
		flag = true
		restTime = interval
		clock.Text = ToTime(restTime)
		clock.Refresh()
	})

	resetBtn := widget.NewButton("重置", func() {
		Stop()
		music = false
		flag = false
		restTime = interval
		clock.Text = ToTime(restTime)
		clock.Refresh()
	})

	defaultBtn := widget.NewButton("默认音乐", func() {
		audioName = ""
		empty.Text = "默认"
		empty.Refresh()
		writeConfig(path, audioName)
	})

	timeSelector := widget.NewSelect([]string{"5", "15", "30", "45", "60"}, func(s string) {
		num, _ := strconv.Atoi(s)
		interval = num * 60
	})
	timeSelector.Selected = "45"

	contentBtn := container.NewCenter(container.New(layout.NewVBoxLayout(), container.New(layout.NewHBoxLayout(), selectMusicBtn, startBtn, resetBtn), container.New(layout.NewHBoxLayout(), defaultBtn, timeSelector)))

	content2 := container.New(layout.NewCenterLayout(), clock)

	clock.Text = ToTime(restTime)
	clock.TextSize = 20
	clock.Refresh()

	w.SetContent(container.New(layout.NewVBoxLayout(), content1, content2, container.New(layout.NewCenterLayout(), empty), contentBtn))

	w.Show()

	go func() {
		for range time.Tick(time.Second) {
			if flag == false {
				continue
			}
			if restTime <= 0 && !music {
				music = true
				go func() {
					if audioName != "" {
						var mp3File []byte
						var err error
						if mp3File, err = ioutil.ReadFile(audioName); err != nil {
							log.Fatal(err)
						}
						play(mp3File)
					} else {
						var nowIndex = rand.Intn(len(audio))
						for lastIndex == nowIndex {
							nowIndex = rand.Intn(len(audio))
						}
						lastIndex = nowIndex
						play(audio[nowIndex])
					}
					music = false
				}()
				continue
			}
			if restTime <= 0 {
				continue
			}
			restTime = restTime - 1
			clock.Text = ToTime(restTime)
			clock.Refresh()
		}
	}()
}

func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//楷体:simkai.ttf
		//黑体:simhei.ttf
		if strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func ToTime(rest int) string {
	return strconv.Itoa(rest/60) + ":" + strconv.Itoa(rest%60)
}

func Stop() {
	if music {
		out <- 1
	}
}

func play(file []byte) {
	var err error
	var dec *minimp3.Decoder
	if dec, err = minimp3.NewDecoder(bytes.NewReader(file)); err != nil {
		log.Fatal(err)
	}
	started := dec.Started()
	<-started

	if context == nil {
		if context, err = oto.NewContext(dec.SampleRate, dec.Channels, 2, 1024*10); err != nil {
			log.Fatal(err)
		}

	}
	var player = context.NewPlayer()

	var waitForPlayOver = new(sync.WaitGroup)
	waitForPlayOver.Add(1)

	go func() {
		for {
			var data = make([]byte, 1024*10)
			_, err := dec.Read(data)
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			player.Write(data)
			select {
			case <-out:
				waitForPlayOver.Done()
				runtime.Goexit()
			default:
			}
		}
		waitForPlayOver.Done()
	}()
	waitForPlayOver.Wait()
	<-time.After(time.Second)
	dec.Close()
	player.Close()
}
