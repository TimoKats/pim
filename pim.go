package main

import (
  lib "github.com/TimoKats/pim/commands/lib"
  pim "github.com/TimoKats/pim/commands"

  "errors"
  "time"
  "os"
)

func heartbeat() {
  lib.Warn.Println("Starting the heartbeat for scheduled tasks. Run this in background!")
  for {
    time.Sleep(time.Second)
  }
}

func parseCommand(command []string, process lib.Process, database *lib.Database) error  {
  switch command[1] {
    case "run":
      if len(command) < 3 {
        return errors.New("No command name given. pim run <<name>>.")
      }
      return pim.RunSelected(command[2], process, database)
    case "start":
      pim.RunSchedule(process, database)
      heartbeat()
    case "log":
      lib.Info.Println("oh you want the logs I see!")
      pim.ViewLogs(database)
    }
  return nil
}

func main() {
  if len(os.Args) < 2 {
    lib.Error.Println("Not enough arguments. pim <<run, start, ls>>.")
    return
  }
  process, database, setupErr := pim.SetupYamlFiles()
  if setupErr != nil {
    lib.Error.Println(setupErr)
    return
  }
  parseErr := parseCommand(os.Args, process, &database)
  if parseErr != nil {
    lib.Error.Println(parseErr)
  }
}
