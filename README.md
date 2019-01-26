# llog - Golang Level Logger 

[![Documentation](https://godoc.org/github.com/midstar/llog?status.svg)](https://godoc.org/github.com/midstar/llog)
[![Go Report Card](https://goreportcard.com/badge/github.com/midstar/llog)](https://goreportcard.com/report/github.com/midstar/llog)
[![AppVeyor](https://ci.appveyor.com/api/projects/status/github/midstar/llog?svg=true)](https://ci.appveyor.com/api/projects/status/github/midstar/llog)
[![Coverage Status](https://coveralls.io/repos/github/midstar/llog/badge.svg?branch=master)](https://coveralls.io/github/midstar/llog?branch=master)


Package llog (Level Logger) extends the standard golang log package with:

* configurable log levels
* log file wrapping if configurable size exceeded

llog is using the "standard" logger in the log package. 

Super simple usage and no configuration required except specifying
the log file name (stderr is used if omitted). 

llog is thread safe and can be used from more than one goroutine.

## Install

```bash
go get github.com/midstar/llog
```

## Import

```go
import (
	"github.com/midstar/llog"
)
```

## Example Usage

```go
package main

import (
	"github.com/midstar/llog"
)

func main() {

	// Only write Info level and above to log
	llog.SetLevel(llog.LvlInfo)
	
	// Write log to mylog.txt. If the log exceeds 1024 KB (=1 MB)
	// a backup with name "mylog.txt.1" will be created and 
	// "mylog.txt" will start over.
	llog.SetFile("mylog.txt", 1024)
	
	// Write some log entries
	llog.Info("This is an info entry. Parameter %d", 23)
	llog.Warn("This is a warning entry")
	llog.Trace("This entry will not be in the log since SetLevel is LvlInfo")
}
```

The file "mylog.txt" will look like this:

	2019/01/26 22:57:15 example.go:18: INFO - This is an info entry. Parameter 23
	2019/01/26 22:57:15 example.go:19: WARN - This is a warning entry

## Notes

You can combine the standard log functions with llog to for example set
the output format. However log file wrapping will only be done from
one of the llog log functions. 

## Author and license

This library is written by Joel Midstjärna and is licensed under the MIT License.