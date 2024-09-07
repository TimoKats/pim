package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"
)

func TrimDatabase(database *lib.Database, threshold int) {
  if len(database.Logs) > threshold {
    database.Logs = database.Logs[len(database.Logs) - threshold:]
  }
}

func CleanDatabase(database *lib.Database) error {
  database.Logs = nil
  lib.WriteDataYaml(lib.DATAPATH, *database)
  return nil
}

