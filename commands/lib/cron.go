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

func TestCron(run Run, process Process, database *Database) (*gocron.Job, error) {
  switch {
    case strings.HasPrefix(run.Schedule, "@times;"):
      return Schedule.Every(1).Day().At(run.Schedule[7:]).Do( func () {
        Info.Println("This is a run to test the schedule.")
      })
    case strings.HasPrefix(run.Schedule, "@start"):
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
