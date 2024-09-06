package lib

import (
  "time"
  "os"
)

func DefaultConfigDir() string { // NOTE: Maybe put this in a variable?
  dirname, _ := os.UserHomeDir()
  return dirname + "/.pim/"
}

func LogfileName() string {
  currentTime := time.Now()
  return currentTime.Format("2006-01-02") + ".log"
}

var IDCHARSET string = "abcdefghijklmnopqrstuvwxyz"
var PROCESSPATH string = DefaultConfigDir() + "process.yaml"
var DATAPATH string = DefaultConfigDir() + "data.yaml"
var LOGPATH string = DefaultConfigDir() + LogfileName()
var COLUMNWIDTH int = 20

