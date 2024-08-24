package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

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

func ViewLogs (database *lib.Database) error {
  for _, log := range database.Logs {
    timeString := log.Timestamp.Format(time.RFC822Z)
    name := responsiveWhitespace(log.RunCommand.Name)
    lib.Info.Printf("| %d\t | %s | %s", log.ExitCode, name, timeString)
  }
  return nil
}
