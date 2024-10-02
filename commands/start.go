// Contains functions related to starting the scheduler (pim start). This module uses an
// external cron module that runs different jobs concurrently.
// Update: I create a lock file that contains the current pid (of the start command) that
// is meant to prevent multiple start processes running at the same time and kill the
// process without calling ps -aux > kill ...
// Note to self: improve this docstring

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "strings"
  "errors"
  "time"

  "github.com/go-co-op/gocron"
)

var schedule *gocron.Scheduler

func heartbeat(process lib.Process, database *lib.Database) {
  lib.Warn.Println("Starting the heartbeat for scheduled tasks. Run this in background!")
  for {
    time.Sleep(10 * time.Second)
    lib.TrimDatabase(database, process.MaxLogs)
    if checkpointErr := lib.WriteCheckpoint(process.Runs, schedule); checkpointErr != nil {
      lib.Error.Println(checkpointErr)
    }
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
  checkpoint, checkpointErr := lib.ReadCheckpoint()
  if checkpointErr != nil {
    lib.Error.Println(checkpointErr)
    return
  }
  for _, run := range checkpoint.Runs {
    if run.Next.Before(time.Now()) {
      lib.Info.Printf("Catch up for '%s'", run.Name)
      schedule.RunByTag(run.Name)
    }
  }
}

func selectCron(run lib.Run, process lib.Process, database *lib.Database) (*gocron.Job, error) {
  if strings.HasPrefix(run.Schedule, "@times;") {
    return schedule.Every(1).Day().At(run.Schedule[7:]).Do( func () {
      lib.Info.Printf("Now running '%s'", run.Name)
      lib.RunAndStore(run, database, process, false)
    })
  }
  return schedule.Cron(run.Schedule).Do( func () {
    lib.Info.Printf("Now running '%s'", run.Name)
    lib.RunAndStore(run, database, process, false)
  })
}

func StartCommand(process lib.Process, database *lib.Database) error {
  if setupErr := setupStart(); setupErr != nil {
    return setupErr
  }
  schedule = gocron.NewScheduler(time.Local)
  for _, run := range process.Runs {
    run := run
    cronJob, cronErr := selectCron(run, process, database)
    cronJob.Tag(run.Name)
    if cronErr != nil {
      lib.Error.Printf("Error in '%s'. %v.", run.Name, cronErr)
    }
  }
  schedule.StartAsync()
  catchup()
  heartbeat(process, database)
  return nil
}
