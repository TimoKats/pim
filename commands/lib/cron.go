package lib

import (
  "strconv"
  "strings"
  "time"

  "github.com/go-co-op/gocron"
)

var Schedule *gocron.Scheduler = gocron.NewScheduler(time.Local)
var RunJobMapping = make(map[*gocron.Job]Run)

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
    if checkpointErr := WriteCheckpoint(Schedule, RunJobMapping); checkpointErr != nil {
      Error.Println(checkpointErr)
    }
    if PIMTERMINATE {
      return
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

