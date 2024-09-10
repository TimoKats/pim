// Contains functions related to cleaning/trimming the files that pim creates. First, the
// the function TrimDatabase is run every heartbeat/run and it's behavior is determined
// by the set_max_logs setting in your process.yaml. Next, the CleanDatabase function
// deletes all the redundant log files and previous runs in data.yaml.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "os"
)

func CleanCommand(database *lib.Database) error {
  database.Logs = nil
  writeErr := lib.WriteDataYaml(lib.DATAPATH, *database)
  if writeErr != nil {
    return writeErr
  }
  os.RemoveAll(lib.LOGDIR)
  makeErr := os.Mkdir(lib.LOGDIR, 0755)
  return makeErr
}

