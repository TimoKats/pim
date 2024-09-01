package lib

import (
  "os"
)

func DefaultConfigDir() string {
  dirname, _ := os.UserHomeDir()
  return dirname + "/.pim/"
}

var IDCHARSET string = "abcdefghijklmnopqrstuvwxyz"
var PROCESSPATH string = DefaultConfigDir() + "process.yaml"
var DATAPATH string = DefaultConfigDir() + "data.yaml"
var COLUMNWIDTH int = 20

