// Most complex module. Contains functions related to executing the runs defined in the
// process yaml. There is a function for a selected run (pim run <run-name>) and a
// scheduled run. The function runAndStore is used by both functions.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
)

func RunCommand(command []string, process lib.Process, database *lib.Database) error {
  var selectedRun string
  if len(command) < 3 {
    return errors.New("No command name given. pim run <<name>>.")
  } else {
    selectedRun = command[2]
  }
  for _, run := range process.Runs {
    if run.Name == selectedRun {
      lib.RunAndStore(run, database, process, true)
      return nil
    }
  }
  return errors.New("'" + command[2] + "' not in process yaml.")
}

