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

  "github.com/go-co-op/gocron"
)

func heartbeat(process lib.Process, database *lib.Database) {
  lib.Warn.Println("Starting the heartbeat for scheduled tasks. Run this in background!")
  for {
    time.Sleep(30 * time.Second)
    lib.TrimDatabase(database, process.MaxLogs)
    lib.WriteCheckpoint()
  }
}

func setupStart() error {
  lib.RemoveDanglingLock()
  if lib.LockExists() {
    return errors.New("Pim is already running! Run <<pim stop>> or check lockfile/ps.")
  }
  lib.InitFileLogging()
  lockErr := lib.InitLockFile()
  if lockErr != nil {
    return lockErr
  }
  return nil
}

func catchup() {

}

func StartCommand(process lib.Process, database *lib.Database) error {
  if setupErr := setupStart(); setupErr != nil {
    return setupErr
  }
  cronSchedule := gocron.NewScheduler(time.Local)
  for _, run := range process.Runs {
    run := run
    _, cronErr := cronSchedule.Cron(run.Schedule).Do( func ()  {
       lib.Info.Printf("Now running '%s'", run.Name)
       lib.RunAndStore(run, database, process, false)
    })
    if cronErr != nil {
      lib.Error.Printf("Error in '%s'. Check Yaml.", run.Name)
    }
  }
  cronSchedule.StartAsync()
  heartbeat(process, database)
  return nil
}
