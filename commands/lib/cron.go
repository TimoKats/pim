// Submodule that sets up the cron schedule. Uses external library called gocron. The
// Schedule variable is used as a 'public' variable to add cronjobs to. There are also
// some functions that are not related to this variable. For example, we allow running on
// startup (not a cron). These functions are also here. Finally, we also have functions
// for testing cron rather than executing it. Those functions contain "dummy" in their
// name.
//
// Note, Heartbeat function is there to keep pim running so that it can execute cron.

package lib

import (
  "strconv"
  "strings"
  "time"

  "github.com/go-co-op/gocron"
)

var Schedule *gocron.Scheduler = gocron.NewScheduler(time.Local)

func Catchup() {
  checkpoint, checkpointErr := ReadCheckpoint()
  if checkpointErr != nil {
    Error.Println(checkpointErr)
    return
  }
  for _, run := range checkpoint.Runs {
    if run.Next.Before(time.Now()) && run.Catchup {
      Info.Printf("Catch up for '%s'", run.Name)
      runErr := Schedule.RunByTag(run.Name)
      if runErr != nil { Error.Println(runErr) }
    }
  }
}

func RunsCatchup(runName string) bool {
  checkpoint, checkpointErr := ReadCheckpoint()
  if checkpointErr != nil {
    return false
  }
  for _, run := range checkpoint.Runs {
    if run.Next.Before(time.Now()) && run.Catchup && run.Name == runName {
      return true
    }
  }
  return false
}

func RunOnStart(run Run, process Process, database *Database) {
  delay := strings.Split(run.Schedule, "+")
  if len(delay) > 1 {
    if delayInt, err := strconv.Atoi(delay[1]); err == nil {
      time.Sleep(time.Duration(delayInt) * time.Second)
    }
  }
  Info.Printf("Now running '%s'", run.Name)
  RunAndStore(run, database, process, false)
}

func SelectCron(run Run, process Process, database *Database) (*gocron.Job, error) {
  switch {
    case strings.HasPrefix(run.Schedule, "@times;"):
      return Schedule.Every(1).Day().At(run.Schedule[7:]).Do( func () {
        Info.Printf("Now running '%s'", run.Name)
        RunAndStore(run, database, process, false)
      })
    case strings.HasPrefix(run.Schedule, "@start"):
      go RunOnStart(run, process, database)
      return nil, nil
    default:
      return Schedule.Cron(run.Schedule).Do( func () {
        Info.Printf("Now running '%s'", run.Name)
        RunAndStore(run, database, process, false)
      })
  }
}

func Heartbeat(process Process, database *Database) {
  Warn.Println("Starting the heartbeat for scheduled tasks. Run this in background!")
  for {
    time.Sleep(10 * time.Second)
    TrimDatabase(database, process.MaxLogs)
    if checkpointErr := WriteCheckpoint(process.Runs, Schedule); checkpointErr != nil {
      Error.Println(checkpointErr)
    }
  }
}

func dummyCron(dummySchedule *gocron.Scheduler, run Run, process Process, database *Database) (*gocron.Job, error) {
  switch {
    case strings.HasPrefix(run.Schedule, "@times;"):
      return dummySchedule.Every(1).Day().At(run.Schedule[7:]).Do( func () {
        Info.Println("This is a run to test the schedule.")
      })
    case strings.HasPrefix(run.Schedule, "@start"):
      return nil, nil
    default:
      return dummySchedule.Cron(run.Schedule).Do( func () {
        Info.Printf("Now running '%s'", run.Name)
        RunAndStore(run, database, process, false)
      })
  }
}

func CreateDummySchedule(process Process, database *Database) *gocron.Scheduler {
  dummySchedule := gocron.NewScheduler(time.Local)
  for _, run := range process.Runs {
    run := run
    cronJob, cronErr := dummyCron(dummySchedule, run, process, database)
    if cronJob != nil && cronErr == nil {
      cronJob.Tag(run.Name)
    }
  }
  dummySchedule.StartAsync()
  return dummySchedule
}

