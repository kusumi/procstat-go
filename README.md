procstat-go ([v0.2.3](https://github.com/kusumi/procstat-go/releases/tag/v0.2.3))
========

## About

+ Ncurses based file monitor.

+ Go version of [https://github.com/kusumi/procstat/](https://github.com/kusumi/procstat/).

## Requirements

go 1.18 or above

## Build

    $ make

## Usage

    $ ./procstat-go -h
    Usage: procstat-go [options] /proc/...
      -bg string
            Set background color. Available colors are "black", "blue", "cyan", "green", "magenta", "red", "white", "yellow".
      -c string
            Set column layout. e.g. "-c 123" to make 3 columns with 1,2,3 windows for each
      -d    Enable debug log
      -f    Fold lines when longer than window width
      -fg string
            Set foreground color. Available colors are "black", "blue", "cyan", "green", "magenta", "red", "white", "yellow".
      -h    This option
      -m    Take refresh interval as milli second. e.g. "-t 500 -m" to refresh screen every 500 milli seconds
      -n    Show line number
      -noblink
            Disable blink
      -r    Rotate column layout
      -t int
            Set refresh interval in second. Default is 1. e.g. "-t 5" to refresh screen every 5 seconds (default 1)
      -usedelay
            Add random delay time before each window starts
      -v    Print version and exit
    
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
      CTRL-l - Repaint whole screen
