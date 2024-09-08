package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"

  "github.com/robfig/cron"
)

func runAndStore(run lib.Run, database *lib.Database, process lib.Process, showOutput bool) {
  var output string
  var status int
  if run.Duration != 0 {
    output, status = lib.ExecuteTimedRun(run, showOutput, run.Duration)
  } else {
    output, status = lib.ExecuteRun(run, showOutput)
  }
  storedLog := lib.StoreRun(run, output, status)
  database.Logs = append(database.Logs, storedLog)
  if (!process.OnlyStoreErrors) || (process.OnlyStoreErrors && status != 0) {
    lib.WriteDataYaml(lib.DATAPATH, *database)
  }
}

func RunSelected(selectedRun string, process lib.Process, database *lib.Database) error {
  for _, run := range process.Runs {
    if run.Name == selectedRun {
      runAndStore(run, database, process, true)
      return nil
    }
  }
  return errors.New("Name of selected run not in process yaml.")
}
 
func RunSchedule(process lib.Process, database *lib.Database) {
  lib.InitFileLogging()
  cronSchedule := cron.New()
  for _, run := range process.Runs {
    run := run
    cronErr := cronSchedule.AddFunc(run.Schedule, func() {
      lib.Info.Printf("Now running '%s'", run.Name)
      runAndStore(run, database, process, false)
    })
    if cronErr != nil {
      lib.Error.Printf("Error in '%s'. Check Yaml.", run.Name)
    }
  }
  cronSchedule.Start()
}

