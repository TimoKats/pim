package lib

import (
  "log"
  "os"
)

var (
  Action *log.Logger
  Warn *log.Logger
  Info *log.Logger
  Error *log.Logger
  Fatal *log.Logger
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"

func init() {
  Info = log.New(os.Stdout, Blue + "info:     " + Reset, log.Ltime|log.Lshortfile)
  Warn = log.New(os.Stdout, Yellow + "warning:  " + Reset, log.Ltime|log.Lshortfile)
  Error = log.New(os.Stdout, Red + "error:    " + Reset, log.Ltime|log.Lshortfile)
  Fatal = log.New(os.Stdout, Magenta + "fatal:    " + Reset, log.Ltime|log.Lshortfile)
}

