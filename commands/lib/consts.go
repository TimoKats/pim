package lib

import (
  "os"
)

func defaultConfigDir() string {
  dirname, _ := os.UserHomeDir()
  return dirname + "/.pim/"
}

var PROCESSPATH string = defaultConfigDir() + "process.yaml"
var DATAPATH string = defaultConfigDir() + "data.yaml"

