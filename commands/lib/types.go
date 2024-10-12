// Datatypes used in pim. Relates to yaml files: process.yaml, data.yaml, checkpoint.
//
// A semantic explanation of how the datatypes are used:
// A <process> consists of <run>s. A <run> creates a <log> that we put in the <database>.
// Checkpoints are created on heartbeats and are not related to the other data types.

package lib

import (
  "time"
)

// process

type Process struct {
  Runs []Run `yaml:"process"`
  OnlyStoreErrors bool `yaml:"only_store_errors"`
  MaxLogs int `yaml:"max_logs"`
}

// storage

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
  Command string `yaml:"command"`
  Duration int `yaml:"duration"`
  Catchup bool `yaml:"catchup"`
}

// checkpoints

type RunCheckpoint struct {
  Next time.Time
  Name string
  Catchup bool
}

type Checkpoint struct {
  Updated time.Time
  Runs []RunCheckpoint
}

