package lib

import (
  "errors"
  "time"
  "os"
)

func DefaultConfigDir() (string, error) {
  dirname, dirErr := os.UserHomeDir()
  if dirErr != nil {
    return "", dirErr
  }
  if _, fileErr := os.Stat(dirname + "/.pim"); os.IsNotExist(fileErr) {
    return "", errors.New("Please create config folder at ~/.pim")
  }
  return dirname + "/.pim/", nil
}

func DefaultLogDir(configDir string) (string, error) {
  logDir := configDir + "/logs/"
  if len(configDir) == 0 { // then the config dir doesn't exist...
    return "", nil
  }
  if _, err := os.Stat(logDir); os.IsNotExist(err) {
    dirErr := os.Mkdir(logDir, 0755)
    if dirErr != nil {
      return "", dirErr
    }
  }
  return logDir, nil
}

func DefaultLogPath() string {
  currentTime := time.Now()
  return currentTime.Format("2006-01-02") + ".log"
}

// somewhat anomalous, but it's a const and it's used so it belongs here...
const HELPSTRING = `Usage: pim <<command>>

  Commands:
  - run <<command-name>>: Runs a command by the name defined in your process YAML.
  - start: Starts the cron schedule defined in your process YAML.
  - stop: Stops the cron schedule started by running: pim start.
  - ls: Lists all the commands and their characteristics defined in your process YAML.
  - log <<optional:run-id>>: Show all logs, or a log of a specific run.
  - clean: Clean log files.
  - stat: Show runs/error rates of the commands defined in your YAML.
  `

// these are checked on startup
var CONFIGDIR, CONFIGERR = DefaultConfigDir()
var LOGDIR, LOGERR = DefaultLogDir(CONFIGDIR)

// these are used by other functions
var IDCHARSET string = "abcdefghijklmnopqrstuvwxyz"
var PROCESSPATH string = CONFIGDIR + "process.yaml"
var DATAPATH string = CONFIGDIR + "data.yaml"
var LOCKPATH string = CONFIGDIR + "lockfile"
var CHECKPOINTPATH string = CONFIGDIR + "checkpoint"
var LOGPATH string = LOGDIR + DefaultLogPath()
var COLUMNWIDTH int = 20
var COMMANDS []string = []string{"ls", "run", "start", "stop", "clean", "stat"}

// meta info
var VERSION string = "v0.0.1"
var AUTHOR string = "Timo Kats"
var LICENSE string = "TBD!"
