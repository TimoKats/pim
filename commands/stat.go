// Contains a function that shows some statistics (runs, errors) for each run in your
// current process.yaml. This module can be elaborated quite a bit probably but contains
// enough for me.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"
)

func countErrors(logs []lib.Log, processName string) int {
  errorCount := 0
  for _, log := range logs {
    if log.ExitCode != 0 && (processName == "" || log.RunCommand.Name == processName) {
      errorCount += 1
    }
  }
  return errorCount
}

func countRuns(logs []lib.Log, processName string) int {
  runCount := 0
  for _, log := range logs {
    if processName == "" || log.RunCommand.Name == processName {
      runCount += 1
    }
  }
  return runCount
}

func StatCommand(process lib.Process, database *lib.Database) error {
  totalRuns := countRuns(database.Logs, "")
  totalErrors := countErrors(database.Logs, "")
  for _, run := range process.Runs {
    runs := countRuns(database.Logs, run.Name)
    errors := countErrors(database.Logs, run.Name)
    lib.Info.Printf("%s \n\truns:   %d\n\terrors: %d", run.Name, runs, errors)
  }
  lib.Info.Printf("total \n\truns:   %d\n\terrors: %d", totalRuns, totalErrors)
  return nil
}

