package main

import (
  lib "github.com/TimoKats/pim/commands/lib"
  pim "github.com/TimoKats/pim/commands"

  "errors"
  "os"
)


func parseCommand(command []string, process lib.Process, database *lib.Database) error  {
  switch command[1] {
    case "run":
      return pim.RunCommand(command, process, database)
    case "start":
      pim.StartCommand(process, database)
    case "log":
      return pim.LogCommand(command, database)
    case "ls":
      return pim.ListCommand(process)
    case "clean":
      return pim.CleanCommand(database)
    case "info":
      return pim.InfoCommand()
    case "stat":
      return pim.StatCommand(process, database)
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
  lib.TrimDatabase(&database, process.MaxLogs)
  if setupErr != nil {
    lib.Error.Println(setupErr)
    return
  }
  parseErr := parseCommand(os.Args, process, &database)
  if parseErr != nil {
    lib.Error.Println(parseErr)
  }
}
