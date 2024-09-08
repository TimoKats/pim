package lib

import (
  "time"
)

type Process struct {
  Runs []Run `yaml:"process"`
  OnlyStoreErrors bool `yaml:"only_store_errors"`
  MaxLogs int `yaml:"max_logs"`
}

type Database struct {
  Logs []Log
}

type Log struct {
  Id string
  RunCommand Run
  Output string
  ExitCode int
  Timestamp time.Time
}

type Run struct {
  Name string `yaml:"name"`
  Directory string `yaml:"directory"`
  Schedule string `yaml:"schedule"`
  Command string `yaml:"command"` // NOTE: Parse this to be split on spaces!
  Duration int `yaml: "duration"`
}
