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

var (
	version [3]int = [3]int{0, 2, 3}
)

func getVersionString() string {
	return fmt.Sprintf("%d.%d.%d", version[0], version[1], version[2])
}

func printVersion() {
	fmt.Println(getVersionString())
}

func usage(progname string) {
	fmt.Fprintln(os.Stderr, "Usage: "+progname+" [options] /proc/...")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, `Commands:
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

type watch struct {
	watcher *fsnotify.Watcher
	fmap    map[string]*window
}

var opt = struct {
	layout    []int
	sinterval int
	minterval int
	showlnum  bool
	foldline  bool
	rotatecol bool
	debug     bool
	fgcolor   int16
	bgcolor   int16
	blinkline bool
	usedelay  bool
}{}

func main() {
	progname := path.Base(os.Args[0])

	optc := flag.String("c", "", "Set column layout. "+
		"e.g. \"-c 123\" to make 3 columns with 1,2,3 windows for each")
	optt := flag.Int("t", 1, "Set refresh interval in second. Default is 1. "+
		"e.g. \"-t 5\" to refresh screen every 5 seconds")
	optm := flag.Bool("m", false, "Take refresh interval as milli second. "+
		"e.g. \"-t 500 -m\" to refresh screen every 500 milli seconds")
	optn := flag.Bool("n", false, "Show line number")
	optf := flag.Bool("f", false, "Fold lines when longer than window width")
	optr := flag.Bool("r", false, "Rotate column layout")
	opth := flag.Bool("h", false, "This option")
	optd := flag.Bool("d", false, "Enable debug log")
	optfg := flag.String("fg", "", "Set foreground color. Available colors are "+
		"\"black\", \"blue\", \"cyan\", \"green\", \"magenta\", \"red\", \"white\", \"yellow\".")
	optbg := flag.String("bg", "", "Set background color. Available colors are "+
		"\"black\", \"blue\", \"cyan\", \"green\", \"magenta\", \"red\", \"white\", \"yellow\".")
	optnoblink := flag.Bool("noblink", false, "Disable blink")
	optusedelay := flag.Bool("usedelay", false, "Add random delay time before each window starts")
	optv := flag.Bool("v", false, "Print version and exit")

	flag.Parse()
	args := flag.Args()
	opt.sinterval = *optt
	opt.minterval = 0
	opt.showlnum = *optn
	opt.foldline = *optf
	opt.rotatecol = *optr
	opt.debug = *optd
	opt.fgcolor = stringToColor(*optfg)
	opt.bgcolor = stringToColor(*optbg)
	opt.blinkline = !*optnoblink
	opt.usedelay = *optusedelay

	if *optv {
		printVersion()
		os.Exit(1)
	}

	if *opth {
		usage(progname)
		os.Exit(1)
	}

	if *optm {
		x := opt.sinterval
		opt.sinterval /= 1000
		opt.minterval = x % 1000
	}
	t := getSecond(opt.sinterval) + getMillisecond(opt.minterval)

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

	if _, errno := getTerminalInfo(); errno != 0 {
		fmt.Fprintln(os.Stderr, "Failed to get terminal info,", errno)
		os.Exit(1)
	}

	defer cleanupLock()
	initLock()
	defer cleanupLog()
	if err := initLog(progname); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init log,", err)
		os.Exit(1)
	}
	defer cleanupScreen()
	if err := initScreen(opt.fgcolor, opt.bgcolor); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init screen,", err)
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init fsnotify,", err)
		os.Exit(1)
	}
	defer watcher.Close()
	watch := watch{watcher, make(map[string]*window)}

	dbg(os.Args)
	dbgf("%#v", opt)

	co := container{}
	co.init(args, &watch)
	dbg(watch.fmap)
	defer co.cleanup()

	sigintCh := make(chan int)
	sigwinchCh := make(chan int)
	exitCh := make(chan int)

	var wg sync.WaitGroup

	// signal handler goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGWINCH)
		for {
			select {
			case <-exitCh:
				dbgf("signal=%p exit", ch)
				return
			case s := <-ch:
				dbg("signal,", s)
				switch s {
				case syscall.SIGINT:
					sigintCh <- 1
				case syscall.SIGWINCH:
					sigwinchCh <- 1
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
			case <-exitCh:
				dbgf("co=%p exit", co)
				return
			case <-sigwinchCh:
				co.parseEvent(KEY_RESIZE)
			default:
				if co.parseEvent(readIncoming()) == -1 {
					dbg("quit")
					sigintCh <- 1
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
			case <-exitCh:
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
						w.updateBuffer()
						flashTerminal()
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
		go func(w *window) {
			defer wg.Done()
			w.repaint()
			d := t
			if opt.usedelay {
				d = getMillisecond(rand.Intn(1000))
			}
			for {
				select {
				case <-exitCh:
					dbgf("window=%p exit", w)
					return
				case <-w.sigCh:
					w.repaint()
				case <-time.After(d):
					w.repaint()
				}
				d = t
			}
		}(w)
	}

	<-sigintCh
	close(exitCh)

	wg.Wait()
}
