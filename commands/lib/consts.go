package lib

import (
  "os"
)

func defaultConfigDir() string {
  dirname, _ := os.UserHomeDir()
  return dirname + "/.pim/"
}

var IDCHARSET string = "abcdefghijklmnopqrstuvwxyz"
var PROCESSPATH string = defaultConfigDir() + "process.yaml"
var DATAPATH string = defaultConfigDir() + "data.yaml"
var COLUMNWIDTH int = 20

