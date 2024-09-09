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
  lib.WriteDataYaml(lib.DATAPATH, *database)
}

func CleanDatabase(database *lib.Database) error {
  database.Logs = nil
  lib.WriteDataYaml(lib.DATAPATH, *database)
  os.RemoveAll(lib.LOGDIR)
  os.Mkdir(lib.LOGDIR, 0755)
  return nil
}

