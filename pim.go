package main

import (
  lib "github.com/TimoKats/pim/commands/lib"
  pim "github.com/TimoKats/pim/commands"

  "errors"
  "time"
  "os"
)

func heartbeat(process lib.Process, database *lib.Database) {
  lib.Warn.Println("Starting the heartbeat for scheduled tasks. Run this in background!")
  for {
    time.Sleep(10 * time.Second)
    pim.TrimDatabase(database, process.MaxLogs)
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
      heartbeat(process, database)
    case "log":
      if len(command) < 3 {
        return pim.ViewLogs(database)
      }
      return pim.ViewLog(database, command[2])
    case "ls":
      return pim.ListProcesses(process)
    case "clean":
      return pim.CleanDatabase(database)
    case "info":
      return pim.Info()
    case "stat":
      return pim.GetStatistics(process, database)
    default:
      return errors.New("Command not found.")
    }
  return nil
}

func main() {
  if len(os.Args) < 2 {
    lib.Error.Println("Not enough arguments. pim <<run, start, ls, log, clean>>.")
    return
  }
  if startupErr := pim.CheckStartupErrors(); startupErr != nil {
    lib.Error.Println(startupErr)
    return
  }
  process, database, setupErr := pim.SetupYamlFiles()
  pim.TrimDatabase(&database, process.MaxLogs)
  if setupErr != nil {
    lib.Error.Println(setupErr)
    return
  }
  parseErr := parseCommand(os.Args, process, &database)
  if parseErr != nil {
    lib.Error.Println(parseErr)
  }
}
