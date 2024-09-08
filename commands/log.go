package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
  "time"
)

func ViewLog (database *lib.Database, logId string) error {
  for _, log := range database.Logs {
    if log.Id == logId {
      timeString := log.Timestamp.Format(time.RFC822Z)
      lib.Info.Printf("exit code: %d", log.ExitCode)
      lib.Info.Printf("command: %s", log.RunCommand.Command)
      lib.Info.Printf("timestamp: %s", timeString)
      lib.Info.Printf("directory: %s", log.RunCommand.Directory)
      lib.Info.Printf("schedule: %s", log.RunCommand.Schedule)
      lib.Info.Println("command output:\n---\n")
      lib.Info.Println(log.Output)
      return nil
    }
  }
  return errors.New("Log id not found in data.")
}

func ViewLogs (database *lib.Database) error {
  for _, log := range database.Logs {
    timeString := log.Timestamp.Format(time.RFC822Z)
    id := lib.ResponsiveWhitespace(log.Id)
    name := lib.ResponsiveWhitespace(log.RunCommand.Name)
    lib.Info.Printf("%d | %s | %s | %s", log.ExitCode, name, timeString, id)
  }
  return nil
}

