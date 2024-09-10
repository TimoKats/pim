// Contains functions related to starting the scheduler (pim start). This module uses an
// external cron module that runs different jobs concurrently.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "time"

  "github.com/robfig/cron"
)

func heartbeat(process lib.Process, database *lib.Database) {
  lib.Warn.Println("Starting the heartbeat for scheduled tasks. Run this in background!")
  for {
    time.Sleep(10 * time.Second)
    lib.TrimDatabase(database, process.MaxLogs)
  }
}

func StartCommand(process lib.Process, database *lib.Database) {
  lib.InitFileLogging()
  cronSchedule := cron.New()
  for _, run := range process.Runs {
    run := run
    cronErr := cronSchedule.AddFunc(run.Schedule, func() {
      lib.Info.Printf("Now running '%s'", run.Name)
      lib.RunAndStore(run, database, process, false)
    })
    if cronErr != nil {
      lib.Error.Printf("Error in '%s'. Check Yaml.", run.Name)
    }
  }
  cronSchedule.Start()
  heartbeat(process, database)
}
