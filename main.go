package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
)

func usage(arg string) {
	s := "Usage: " + arg + " [options] /proc/...\n" +
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
  CTRL-l - Repaint whole screen`
	fmt.Fprintln(os.Stderr, s)
}

func signal_handler(sigint_ch, sigwinch_ch chan<- bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGWINCH)

	for signal := range c {
		switch signal {
		case syscall.SIGINT:
			dbg("SIGINT")
			sigint_ch <- true
		case syscall.SIGWINCH:
			dbg("SIGWINCH")
			sigwinch_ch <- true
		}
	}
}

type Option struct {
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
}

var opt Option

func main() {
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
		usage(os.Args[0])
		return
	}

	if *optm {
		x := opt.sinterval
		opt.sinterval /= 1000
		opt.minterval = x % 1000
	}
	t := time.Duration(opt.sinterval)*time.Second +
		time.Duration(opt.minterval)*time.Millisecond

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
		return
	}

	defer CleanupLock()
	InitLock()
	defer CleanupLog()
	if err := InitLog(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init log,", err)
		return
	}
	defer CleanupScreen()
	if err := InitScreen(opt.fgcolor, opt.bgcolor); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init screen,", err)
		return
	}

	dbg(os.Args)
	dbgf("%#v", opt)

	co := new(Container)
	co.Init(args)

	sigint_ch := make(chan bool)
	sigwinch_ch := make(chan bool)
	c_ch := make(chan bool)
	w_ch := make(chan bool)

	go signal_handler(sigint_ch, sigwinch_ch)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(co *Container, exit_ch <-chan bool) {
		defer wg.Done()
		for {
			select {
			case <-exit_ch:
				dbg("exit")
				return
			case <-sigwinch_ch:
				dbg("signal/winch")
				co.ParseEvent(KEY_RESIZE)
			default:
				if co.ParseEvent(ReadIncoming()) == -1 {
					dbg("quit")
					sigint_ch <- true
				}
			}
		}
	}(co, c_ch)

	for _, w := range co.v {
		wg.Add(1)
		go func(w *Window, exit_ch <-chan bool) {
			defer wg.Done()
			w.Repaint()
			d := t
			if opt.usedelay {
				r := rand.Intn(1000)
				d = time.Duration(r) * time.Millisecond
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
		}(w, w_ch)
	}

	<-sigint_ch
	close(c_ch)
	close(w_ch)

	wg.Wait()
}
