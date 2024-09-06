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

func ResponsiveWhitespace(text string) string {
  if len(text) > COLUMNWIDTH {
    text = text[:COLUMNWIDTH]
  }
  spaces := COLUMNWIDTH - len(text)
  for i := 0; i < spaces; i++ {
    text += " "
  }
  return text
}

func InitFileLogging() {
  logFile, err := os.OpenFile(LOGPATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
  if err != nil {
    log.Fatal(err)
  }
  Info = log.New(logFile, "info:  ", log.Ltime)
  Warn = log.New(logFile, "warn:  ", log.Ltime)
  Error = log.New(logFile, "error: ", log.Ltime)
  Fatal = log.New(logFile, "fatal: ", log.Ltime)
}

func init() {
  Info = log.New(os.Stdout, Blue + "info:    " + Reset, log.Ltime)
  Warn = log.New(os.Stdout, Yellow + "warning: " + Reset, log.Ltime)
  Error = log.New(os.Stdout, Red + "error:   " + Reset, log.Ltime)
  Fatal = log.New(os.Stdout, Magenta + "fatal:   " + Reset, log.Ltime)
}

