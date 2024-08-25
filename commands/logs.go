package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
  "time"
)

func responsiveWhitespace(text string) string {
  if len(text) > lib.COLUMNWIDTH {
    text = text[:lib.COLUMNWIDTH]
  }
  spaces := lib.COLUMNWIDTH - len(text)
  for i := 0; i < spaces; i++ {
    text += " "
  }
  return text
}

func ViewLog (database *lib.Database, logId string) error {
  for _, log := range database.Logs {
    if log.Id == logId {
      timeString := log.Timestamp.Format(time.RFC822Z)
      cmd := responsiveWhitespace(log.RunCommand.Command)
      dir := responsiveWhitespace(log.RunCommand.Directory)
      schedule := responsiveWhitespace(log.RunCommand.Schedule)
      lib.Info.Printf("%d | %s | %s | %s | %s ", log.ExitCode, cmd, timeString, dir, schedule)
      lib.Info.Println(log.Output)
      return nil
    }
  }
  return errors.New("Log id not found in data.")
}

func ViewLogs (database *lib.Database) error {
  for _, log := range database.Logs {
    timeString := log.Timestamp.Format(time.RFC822Z)
    id := responsiveWhitespace(log.Id)
    name := responsiveWhitespace(log.RunCommand.Name)
    lib.Info.Printf("%d | %s | %s | %s", log.ExitCode, name, timeString, id)
  }
  return nil
}
