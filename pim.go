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
      return pim.StartCommand(process, database)
    case "stop":
      return pim.StopCommand()
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
}

func main() {
  if startupErr := pim.CheckStartupErrors(); startupErr != nil {
    lib.Error.Println(startupErr)
    return
  }
  if len(os.Args) < 2 {
    lib.Error.Println(lib.HELPSTRING)
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
