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
  if len(configDir) == 0 { // then the config dir doesn't exist...
    return "", nil
  }
  if _, err := os.Stat(logDir); os.IsNotExist(err) {
    dirErr := os.Mkdir(logDir + "logs/", 0755)
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

// these are checked on startup
var CONFIGDIR, CONFIGERR = DefaultConfigDir()
var LOGDIR, LOGERR = DefaultLogDir(CONFIGDIR)

// these are used by other functions
var IDCHARSET string = "abcdefghijklmnopqrstuvwxyz"
var PROCESSPATH string = CONFIGDIR + "process.yaml"
var DATAPATH string = CONFIGDIR + "data.yaml"
var LOGPATH string = LOGDIR + DefaultLogPath()
var COLUMNWIDTH int = 20

