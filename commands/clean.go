// Contains functions related to cleaning/trimming the files that pim creates. First, the
// the function TrimDatabase is run every heartbeat/run and it's behavior is determined
// by the set_max_logs setting in your process.yaml. Next, the CleanDatabase function
// deletes all the redundant log files and previous runs in data.yaml.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "os"
)

func TrimDatabase(database *lib.Database, threshold int) {
  if len(database.Logs) > threshold && threshold != 0 {
    database.Logs = database.Logs[len(database.Logs) - threshold:]
  }
  writeErr := lib.WriteDataYaml(lib.DATAPATH, *database)
  if writeErr != nil {
    lib.Warn.Println("Wasn't able to trim database. Continuing operations...")
  }
}

func CleanDatabase(database *lib.Database) error {
  database.Logs = nil
  writeErr := lib.WriteDataYaml(lib.DATAPATH, *database)
  if writeErr != nil {
    return writeErr
  }
  os.RemoveAll(lib.LOGDIR)
  makeErr := os.Mkdir(lib.LOGDIR, 0755)
  return makeErr
}

