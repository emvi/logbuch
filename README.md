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

    // logging cannot be disabled for errors except you use the DiscardFormatter
    logbuch.SetLevel(logbuch.LevelInfo)
    logbuch.Debug("Don't log this anymore!")

    // create your own logger
    l := logbuch.NewLogger(os.Stdout, os.Stderr)
    l.Debug("Just like the default logger...")
    l.SetFormatter(logbuch.NewDiscardFormatter())
    l.Error("This error will be dropped!")
    
    // or to panic...
    l.Fatal("We are going down! Error code: %d", 123)
}
```

## Formatters

To use formatters you can either implement your own or use one provided by logbuch. There are three kind of formatters provided right now:

### StandardFormatter

This is the default. The log output looks like this:

```
2019-09-19T17:39:02.4326139+02:00 [DEBUG] This is a debug message.
2019-09-19T17:39:02.4326139+02:00 [INFO ] Hello World!
2019-09-19T17:39:02.4326139+02:00 [WARN ] Some formatted message 123.
2019-09-19T17:39:02.4326139+02:00 [ERROR] An error occurred: 123
```

## FieldFormatter

The FieldFormatter prints the log parameters in a structured way. To have a nice logging output, use the `logbuch.Fields` type together with this:

```
formatter := logbuch.NewFieldFormatter(logbuch.StandardTimeFormat, "\t\t\t")
logbuch.SetFormatter(formatter)
logbuch.Debug("Debug message", logbuch.Fields{"some": "value", "code": 123})
```

The log output looks like this:

```
2019-09-19T17:45:26.6635897+02:00 [DEBUG] Debug message				 some=value code=123
```

### DiscardFormatter

The DiscardFormatter simply drops all log messages (including errors) and can be used to do just that.

## Contribute

[See CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT
