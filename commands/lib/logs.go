// Sets up logging format for pim. Note, if pim start is used, then switch to file logging.
// This prevents the user being interrupted by writing to standard output when running in
// the background.

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
  Info = log.New(os.Stdout, "", 0)
  Warn = log.New(os.Stdout, Yellow + "warn:  " + Reset, 0)
  Error = log.New(os.Stdout, Red + "error: " + Reset, 0)
  Fatal = log.New(os.Stdout, Magenta + "fatal: " + Reset, 0)
}

