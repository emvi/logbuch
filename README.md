# Logbuch

[![GoDoc](https://godoc.org/github.com/emvi/logbuch?status.svg)](https://godoc.org/github.com/emvi/logbuch)
[![CircleCI](https://circleci.com/gh/emvi/logbuch.svg?style=svg)](https://circleci.com/gh/emvi/logbuch)
[![Go Report Card](https://goreportcard.com/badge/github.com/emvi/logbuch)](https://goreportcard.com/report/github.com/emvi/logbuch)

Simple Go logging library with support for different output channels (io.Writer) for each log level. A formatter can be provided to change the log output formatting.

## Installation

To install logbuch, run go get within your project:

```
go get github.com/emvi/logbuch
```

## Usage

Here is a quick example on how to use the basic functionality of logbuch:

```
import (
    "os"
    "github.com/emvi/logbuch"
)

func main() {
    // use the default logger (logging to stdout and stderr)
    logbuch.Debug("Hello %s!", "World")
    logbuch.Info("Info")
    logbuch.Warn("Warning")
    logbuch.Error("Error")

    logbuch.SetLevel(logbuch.LevelInfo)
    logbuch.Debug("Don't log this anymore!")

    // create your own logger
    l := logbuch.NewLogger(os.Stdout, os.Stderr)
    l.Debug("Just like the default logger...")
    l.SetFormatter(logbuch.NewDiscardFormatter())
    l.Error("This error will be dropped!")
}
```

## Contribute

[See CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT
