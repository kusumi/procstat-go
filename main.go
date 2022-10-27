package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"

	"github.com/fsnotify/fsnotify"
)

func usage(progname string) {
	fmt.Fprintln(os.Stderr, "Usage: "+progname+" [options] /proc/...\n"+
		`Options:
  -c <arg> - Set column layout. e.g. "-c 123" to make 3 columns with 1,2,3 windows for each
  -t <arg> - Set refresh interval in second. Default is 1. e.g. "-t 5" to refresh screen every 5 seconds
  -m - Take refresh interval as milli second. e.g. "-t 500 -m" to refresh screen every 500 milli seconds
  -n - Show line number
  -f - Fold lines when longer than window width
  -r - Rotate column layout
  -h - This option
  --fg <arg> - Set foreground color. Available colors are "black", "blue", "cyan", "green", "magenta", "red", "white", "yellow".
  --bg <arg> - Set background color. Available colors are "black", "blue", "cyan", "green", "magenta", "red", "white", "yellow".
  --noblink - Disable blink
  --usedelay - Add random delay time before each window starts

Commands:
  0 - Set current position to the first line of the buffer
  $ - Set current position to the last line of the buffer
  k|UP - Scroll upward
  j|DOWN - Scroll downward
  h|LEFT - Select next window
  l|RIGHT - Select previous window
  CTRL-b - Scroll one page upward
  CTRL-u - Scroll half page upward
  CTRL-f - Scroll one page downward
  CTRL-d - Scroll half page downward
  CTRL-l - Repaint whole screen`)
}

type Watch struct {
	watcher *fsnotify.Watcher
	fmap    map[string]*Window
}

var opt = struct {
	layout    []int
	sinterval int
	minterval int
	showlnum  bool
	foldline  bool
	rotatecol bool
	usage     bool
	debug     bool
	fgcolor   int16
	bgcolor   int16
	blinkline bool
	usedelay  bool
}{}

func main() {
	progname := path.Base(os.Args[0])

	optc := flag.String("c", "", "")
	optt := flag.Int("t", 1, "")
	optm := flag.Bool("m", false, "")
	optn := flag.Bool("n", false, "")
	optf := flag.Bool("f", false, "")
	optr := flag.Bool("r", false, "")
	opth := flag.Bool("h", false, "")
	optd := flag.Bool("d", false, "")
	optfg := flag.String("fg", "", "")
	optbg := flag.String("bg", "", "")
	optnoblink := flag.Bool("noblink", false, "")
	optusedelay := flag.Bool("usedelay", false, "")

	flag.Parse()
	args := flag.Args()
	opt.sinterval = *optt
	opt.minterval = 0
	opt.showlnum = *optn
	opt.foldline = *optf
	opt.rotatecol = *optr
	opt.usage = *opth
	opt.debug = *optd
	opt.fgcolor = StringToColor(*optfg)
	opt.bgcolor = StringToColor(*optbg)
	opt.blinkline = !*optnoblink
	opt.usedelay = *optusedelay

	if opt.usage {
		usage(progname)
		os.Exit(1)
	}

	if *optm {
		x := opt.sinterval
		opt.sinterval /= 1000
		opt.minterval = x % 1000
	}
	t := GetSecond(opt.sinterval) + GetMillisecond(opt.minterval)

	if *optc == "" {
		*optc = strings.Repeat("1", len(args))
		if *optc == "" {
			*optc = "1"
		}
	}

	opt.layout = make([]int, 0)
	for _, x := range *optc {
		x = unicode.ToUpper(x)
		if x > '0' && x <= '9' {
			opt.layout = append(opt.layout, int(x-'0'))
		} else if x >= 'A' && x <= 'F' {
			opt.layout = append(opt.layout, int(x-'A'+10))
		} else {
			opt.layout = append(opt.layout, -1)
		}
	}

	if _, errno := GetTerminalInfo(); errno != 0 {
		fmt.Fprintln(os.Stderr, "Failed to get terminal info,", errno)
		os.Exit(1)
	}

	defer CleanupLock()
	InitLock()
	defer CleanupLog()
	if err := InitLog(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init log,", err)
		os.Exit(1)
	}
	defer CleanupScreen()
	if err := InitScreen(opt.fgcolor, opt.bgcolor); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init screen,", err)
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init fsnotify,", err)
		os.Exit(1)
	}
	defer watcher.Close()
	watch := Watch{watcher, make(map[string]*Window)}

	dbg(os.Args)
	dbgf("%#v", opt)

	co := Container{}
	co.Init(args, &watch)
	dbg(watch.fmap)

	sigint_ch := make(chan int)
	sigwinch_ch := make(chan int)
	exit_ch := make(chan int)

	var wg sync.WaitGroup

	// signal handler goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGWINCH)
		for {
			select {
			case <-exit_ch:
				dbgf("signal=%p exit", ch)
				return
			case s := <-ch:
				dbg("signal,", s)
				switch s {
				case syscall.SIGINT:
					sigint_ch <- 1
				case syscall.SIGWINCH:
					sigwinch_ch <- 1
				}
			}
		}
	}()

	// container goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-exit_ch:
				dbgf("co=%p exit", co)
				return
			case <-sigwinch_ch:
				co.ParseEvent(KEY_RESIZE)
			default:
				if co.ParseEvent(ReadIncoming()) == -1 {
					dbg("quit")
					sigint_ch <- 1
				}
			}
		}
	}()

	// watch goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-exit_ch:
				dbgf("watch=%p exit", watch)
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				dbg(event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					abs, _ := filepath.Abs(event.Name)
					if w, ok := watch.fmap[abs]; ok {
						w.UpdateBuffer()
					} else {
						dbg("No such key", abs)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				dbgf("watch=%p error=%s", watch, err)
			}
		}
	}()

	// window goroutines
	for _, w := range co.v {
		wg.Add(1)
		go func(w *Window) {
			defer wg.Done()
			w.Repaint()
			d := t
			if opt.usedelay {
				d = GetMillisecond(rand.Intn(1000))
			}
			for {
				select {
				case <-exit_ch:
					dbgf("window=%p exit", w)
					return
				case <-w.sig_ch:
					w.Repaint()
				case <-time.After(d):
					w.Repaint()
				}
				d = t
			}
		}(w)
	}

	<-sigint_ch
	close(exit_ch)

	wg.Wait()
}
