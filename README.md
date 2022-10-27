procstat-go ([v0.2.1](https://github.com/kusumi/procstat-go/releases/tag/v0.2.1))
========

## About

+ Ncurses based file monitor.

+ Go version of [https://github.com/kusumi/procstat/](https://github.com/kusumi/procstat/).

![procstat](https://a.fsdn.com/con/app/proj/procfsv/screenshots/318601.jpg/max/max/1)

## Requirements

+ go 1.18 or above

+ github.com/rthornton128/goncurses

+ github.com/fsnotify/fsnotify

## Build

    $ make

or

    $ gmake

## Usage

    $ ./procstat-go -h
    Usage: procstat-go [options] /proc/...
    Options:
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
      CTRL-l - Repaint whole screen

## Resource

+ Upstream [https://sourceforge.net/projects/procfsv/](https://sourceforge.net/projects/procfsv/)

+ Repository [https://github.com/kusumi/procstat-go/](https://github.com/kusumi/procstat-go/)
