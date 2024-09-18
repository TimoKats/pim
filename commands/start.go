// Contains functions related to starting the scheduler (pim start). This module uses an
// external cron module that runs different jobs concurrently.
// Update: I create a lock file that contains the current pid (of the start command) that
// is meant to prevent multiple start processes running at the same time and kill the
// process without calling ps -aux > kill ...
// Note to self: improve this docstring

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
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

func StartCommand(process lib.Process, database *lib.Database) error {
  if lib.LockExists() {
    return errors.New("Pim is already running! Run <<pim stop>> or check lockfile.")
  }
  lib.InitFileLogging()
  lockErr := lib.InitLockFile()
  if lockErr != nil {
    return lockErr
  }
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
  return nil
}
