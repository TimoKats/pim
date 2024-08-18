package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"

  "github.com/robfig/cron"
)

func runRun (run lib.Run, database *lib.Database) {
  lib.Info.Printf("Now running '%s'", run.Name)
  output, status := lib.ExecuteRun(run, false)
  storedLog := lib.StoreRun(run, output, status)
  database.Logs = append(database.Logs, storedLog)
  lib.WriteDataYaml("tests/data.yaml", *database)
} // NOTE: This can be moved to the run schedule function...

func RunSelected (selectedRun string, process lib.Process, database *lib.Database) error {
  for _, run := range process.Runs {
    if run.Name == selectedRun {
      output, status := lib.ExecuteRun(run, true)
      storedLog := lib.StoreRun(run, output, status)
      database.Logs = append(database.Logs, storedLog)
      lib.WriteDataYaml("tests/data.yaml", *database)
      return nil
    }
  }
  return errors.New("Name of selected run not in process yaml.")
}
  
func RunSchedule(process lib.Process, database *lib.Database) {
  cronSchedule := cron.New()
  for _, run := range process.Runs {
    run := run
    cronErr := cronSchedule.AddFunc(run.Schedule, func() { runRun(run, database) })
    if cronErr != nil {
      lib.Error.Printf("Error in '%s'. Check Yaml.", run.Name)
    }
  }
  cronSchedule.Start()
}

