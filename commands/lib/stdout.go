// Not sure what to call this module. It helps with formatting tables in standard output
// for a number of functions that I use elsewhere. pim log and pim ls use this.

package lib

import (
  "github.com/go-co-op/gocron"

  "strings"
  "errors"
  "time"
)

func ResponsiveWhitespace(text string) string {
  if len(text) > COLUMNWIDTH {
    text = text[:COLUMNWIDTH]
  }
  spaces := COLUMNWIDTH - len(text)
  for i := 0; i < spaces; i++ {
    text += " "
  }
  return text
}

func ViewLog (database *Database, logId string) error {
  for _, log := range database.Logs {
    if log.Id == logId {
      timeString := log.Timestamp.Format(time.RFC822Z)
      Info.Printf("exit code: %d", log.ExitCode)
      Info.Printf("command: %s", log.RunCommand.Command)
      Info.Printf("timestamp: %s", timeString)
      Info.Printf("directory: %s", log.RunCommand.Directory)
      Info.Printf("schedule: %s", log.RunCommand.Schedule)
      Info.Println("command output:\n---")
      Info.Println(log.Output)
      return nil
    }
  }
  return errors.New("Log id not found in data.")
}

func ViewLogs (database *Database) error {
  for _, log := range database.Logs {
    timeString := log.Timestamp.Format(time.RFC822Z)
    id := ResponsiveWhitespace(log.Id)
    name := ResponsiveWhitespace(log.RunCommand.Name)
    Info.Printf("%d | %s | %s | %s", log.ExitCode, name, timeString, id)
  }
  return nil
}

func ViewNextRun (schedule *gocron.Scheduler, run Run) (string, bool) {
  var nextRun string
  var runsCatchup bool
  cronJob, cronErr := schedule.FindJobsByTag(run.Name)
  if cronErr == nil && cronJob != nil {
    nextRun = ResponsiveWhitespace(cronJob[0].NextRun().Format(time.RFC822Z))
    runsCatchup = RunsCatchup(run.Name)
  } else {
    nextRun = ResponsiveWhitespace("No cron schedule.")
    runsCatchup = strings.HasPrefix(run.Schedule, "@start")
  }
  return nextRun, runsCatchup
}

func ViewListHeader() string {
  name := ResponsiveWhitespace("Name")
  cronString := ResponsiveWhitespace("Cron string")
  command := ResponsiveWhitespace("Command")
  duration := ResponsiveWhitespace("Max duration")
  nextRun := ResponsiveWhitespace("Next Run")
  runStart := ResponsiveWhitespace("Runs on start")
  return strings.Join(
    []string{name, cronString, command, duration, nextRun, runStart},
    " | ",
  )
}

