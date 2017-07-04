package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	owm "github.com/briandowns/openweathermap"
	proc "github.com/c9s/goprocinfo/linux"
	"github.com/fhs/gompd/mpd"
	"github.com/fsnotify/fsnotify"
	"github.com/mdlayher/lmsensors"
)

// Bar ...
type Bar struct {
	Done chan bool

	// TODO: I don't really want this here, but I also don't want two
	// XUtils/roots.
	xu   *xgbutil.XUtil
	root *xwindow.Window

	clock        string
	memory       string
	music        string
	musicState   bool
	temperature  string
	todo         string
	weather      string
	weatherState int
	window       string
	workspace    string
}

func newBar() (*Bar, error) {
	b := new(Bar)
	var err error

	b.Done = make(chan bool)

	b.xu, err = xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	b.root = xwindow.New(b.xu, b.xu.RootWin())
	b.root.Listen(xproto.EventMaskPropertyChange)

	// Set default values.
	b.clock = "--:-- --"
	b.memory = "-- MB"
	b.music = "- - -"
	b.musicState = false
	b.temperature = "-- °C"
	b.todo = "0"
	b.weather = "-- °C"
	b.weatherState = 0
	b.window = "-"
	b.workspace = "---"

	return b, nil
}

func (b *Bar) clockFun() {
	init := true
	for {
		if !init {
			time.Sleep(20 * time.Second)
		}
		init = false

		t := time.Now()

		b.clock = t.Format("03:04 PM")
		b.Done <- true
	}
}

func (b *Bar) memoryFun() {
	init := true
	for {
		if !init {
			time.Sleep(20 * time.Second)
		}
		init = false

		mem, err := proc.ReadMemInfo("/proc/meminfo")
		if err != nil {
			log.Print(err)
			continue
		}

		b.memory = strconv.Itoa(int(mem.Active)/1000) + " MB"
		b.Done <- true
	}
}

func (b *Bar) musicFun() {
	watcher, err := mpd.NewWatcher("tcp", ":6600", "", "player")
	if err != nil {
		log.Fatal(err)
	}

	var conn *mpd.Client
	init := true
	for {
		if !init {
			conn.Close()
			<-watcher.Event
		}
		init = false

		conn, err = mpd.Dial("tcp", ":6600")
		if err != nil {
			log.Print(err)
			continue
		}

		status, err := conn.Status()
		if err != nil {
			log.Print(err)
			continue
		}
		cur, err := conn.CurrentSong()
		if err != nil {
			log.Print(err)
			continue
		}

		b.music = strings.ToUpper(cur["Artist"] + " - " +
			cur["Title"])
		b.musicState = status["state"] == "pause"
		b.Done <- true
	}
}

func (b *Bar) temperatureFun() {
	device, err := lmsensors.New().Scan()
	if err != nil {
		log.Fatal(err)
	}
	sensor := device[0].Sensors[0]

	init := true
	var Otemp float64
	for {
		if !init {
			time.Sleep(20 * time.Second)
		}
		init = false

		temp := sensor.(*lmsensors.TemperatureSensor).Input
		if Otemp == temp {
			continue
		}
		Otemp = temp

		b.temperature = strconv.FormatFloat(temp, 'f', 0, 64) + " °C"
		b.Done <- true
	}
}

func (b *Bar) todoFun() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	if err := watcher.Add("/home/onodera/todo"); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("/home/onodera/todo")
	if err != nil {
		log.Fatal(err)
	}

	init := true
	for {
		if !init {
			ev := <-watcher.Events
			if ev.Op&fsnotify.Write != fsnotify.Write {
				continue
			}
		}
		init = false

		s := bufio.NewScanner(file)
		s.Split(bufio.ScanLines)
		var c int
		for s.Scan() {
			c++
		}
		if _, err := file.Seek(0, 0); err != nil {
			log.Print(err)
			continue
		}

		b.todo = strconv.Itoa(c)
		b.Done <- true
	}
}

func (b *Bar) weatherFun() {
	w, err := owm.NewCurrent("C", "en")
	if err != nil {
		log.Fatalln(err)
	}

	init := true
	for {
		if !init {
			time.Sleep(200 * time.Second)
		}
		init = false

		if err := w.CurrentByID(2758106); err != nil {
			log.Print(err)
			continue
		}

		var state int
		switch w.Weather[0].Icon[0:2] {
		case "01":
			state = 0
		case "02":
			state = 1
		case "03":
			state = 2
		case "04":
			state = 3
		case "09":
			state = 4
		case "10":
			state = 5
		case "11":
			state = 6
		case "13":
			state = 7
		case "50":
			state = 8
		}

		b.weather = strconv.FormatFloat(w.Main.Temp, 'f', 0, 64) +
			" °C"
		b.weatherState = state
		b.Done <- true
	}
}

func (b *Bar) windowFun() {
	init := true
	var Owin string
	for {
		if !init {
			ev, xgbErr := b.xu.Conn().WaitForEvent()
			if xgbErr != nil {
				log.Print(xgbErr)
				continue
			}

			atom, err := xprop.Atm(b.xu, "_NET_ACTIVE_WINDOW")
			if ev.(xproto.PropertyNotifyEvent).Atom != atom {
				continue
			}
			if err != nil {
				log.Print(err)
				continue
			}
		}
		init = false

		id, err := ewmh.ActiveWindowGet(b.xu)
		if err != nil {
			log.Print(err)
			continue
		}

		win, err := ewmh.WmNameGet(b.xu, id)
		if err != nil {
			log.Print(err)
			continue
		}
		if Owin == win {
			continue
		}
		Owin = win

		b.window = strings.ToUpper(win)
		b.Done <- true
	}
}

func (b *Bar) workspaceFun() {
	init := true
	var Owsp string
	for {
		if !init {
			ev, xgbErr := b.xu.Conn().WaitForEvent()
			if xgbErr != nil {
				log.Print(xgbErr)
				continue
			}

			atom, err := xprop.Atm(b.xu, "WINDOWCHEF_ACTIVE_GROUPS")
			if ev.(xproto.PropertyNotifyEvent).Atom != atom {
				continue
			}
			if err != nil {
				log.Print(err)
				continue
			}
		}
		init = false

		n, err := ewmh.CurrentDesktopGet(b.xu)
		if err != nil {
			log.Print(err)
			continue
		}

		var wsp string
		switch n {
		case 0:
			wsp = "WWW"
		case 1:
			wsp = "IRC"
		case 2:
			wsp = "SRC"
		}
		if Owsp == wsp {
			continue
		}
		Owsp = wsp

		b.workspace = wsp
		b.Done <- true
	}
}
