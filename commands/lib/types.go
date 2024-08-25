package lib

import (
  "time"
)

type Process struct {
  Runs []Run `yaml:"process"`
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
}
