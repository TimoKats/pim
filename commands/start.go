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

func heartbeat(process lib.Process, database *lib.Database, schedule *gocron.Scheduler) {
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

// func catchup(schedule *gocron.Scheduler) {
//   _, checkpointErr := lib.ReadCheckpoint()
//   if checkpointErr != nil {
//     lib.Error.Println(checkpointErr)
//     return
//   }
//   for _, job := range schedule.Jobs() {
//     lib.Info.Println(job.NextRun())
//   }
// }

func selectCron(
  run lib.Run, process lib.Process, database *lib.Database,
  schedule *gocron.Scheduler) (*gocron.Job, error) { // NOTE: schedule can be global var
  if strings.HasPrefix(run.Schedule, "@times;") {
    return schedule.Every(1).Month(2).At(run.Schedule[7:]).Do( func () {
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
  schedule := gocron.NewScheduler(time.Local)
  for _, run := range process.Runs {
    run := run
    _, cronErr := selectCron(run, process, database, schedule)
    if cronErr != nil {
      lib.Error.Printf("Error in '%s'. %v.", run.Name, cronErr)
    }
  }
  schedule.StartAsync()
  heartbeat(process, database, schedule)
  return nil
}
