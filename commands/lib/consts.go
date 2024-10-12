// This submodule contains all constant variables used in pim. Some are based on the home
// directory of the user. Hence, there are also some functions to derive that (no xgd).
// The const variables are: paths, helpstrings, and configuration settings. Note, use all
// caps format when defining consts.

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

const ASCIILOGO = `
     _         | Description: Pim is a task orchestration tool.
 ___|_|_____   | Version: %s
| . | |     |  | Author: %s
|  _|_|_|_|_|  | License: %s
|_|            | Source code: github.com/TimoKats/pim

`

const ABSTRACT = `
  Abstract:
  Pim is a process orchestration tool (i.e. it can shedule/call commands). To get started,
  setup a <<~/pim>> folder in your home directory. Here, you can write a <<process.yaml>>
  that contains your processes/commands and their schedule(s). Thereafter, you can start
  using Pim. For more information on how to setup your process.yaml, please visit the
  documentation on GitHub.`

const HELPSTRING = `
  Usage:
  - pim <<command>>

  Commands:
  - run <<command-name>>: Runs a command by the name defined in your process YAML.
  - start: Starts the cron schedule defined in your process YAML.
  - stop: Stops the cron schedule started by running: pim start.
  - ls: Lists all the commands and their characteristics defined in your process YAML.
  - log <<optional:run-id>>: Show all logs, or a log of a specific run.
  - clean: Clean log files.

  Flags:
  - info/i: Outputs some information about this Pim installation.
  - help/h: Well...if you see this message you probably typed this...
  - version/v: Shows version of this Pim installation.
  - license/l: Shows the license of this Pim installation.`

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
var COLUMNWIDTH int = 25
var COMMANDS []string = []string{"ls", "run", "start", "stop", "clean", "stat"}

// meta info
var VERSION string = "v0.0.1"
var AUTHOR string = "Timo Kats"
var LICENSE string = "The GNU General Public License v3.0"
