package main

import (
  lib "github.com/TimoKats/pim/commands/lib"
  pim "github.com/TimoKats/pim/commands"

  "testing"
)

var database lib.Database
var process = lib.Process{
  Runs:[]lib.Run{
    lib.Run{
      Name: "linux-command",
      Schedule: "@hourly",
      Command: "echo hello world",
    },
    lib.Run{
      Name: "sleepy-command",
      Schedule: "@hourly",
      Command: "sleep 50",
      Duration: 10,
    },
  },
}


func TestLs(t *testing.T) {
  cmdErr := pim.ListCommand(process)
  if cmdErr != nil {
    t.Errorf("Error in ls command: %v", cmdErr)
  }
}

func TestNonTimedRun(t *testing.T) {
  command := []string{"one", "sleep", "linux-command"}
  cmdErr := pim.RunCommand(command, process, &database)
  if cmdErr != nil {
    t.Errorf("Error in run command: %v", cmdErr)
  }
}

func TestTimedRun(t *testing.T) {
  command := []string{"one", "sleep", "sleepy-command"}
  cmdErr := pim.RunCommand(command, process, &database)
  if cmdErr != nil {
    t.Errorf("Error in run command: %v", cmdErr)
  }
}

func TestLog(t *testing.T) {
  command := []string{"one", "two"}
  cmdErr := pim.LogCommand(command, &database)
  if cmdErr != nil {
    t.Errorf("Error in log command: %v", cmdErr)
  }
}

func TestInfo(t *testing.T) {
  cmdErr := pim.InfoCommand()
  if cmdErr != nil {
    t.Errorf("Error in info command: %v", cmdErr)
  }
}

func TestStat(t *testing.T) {
  cmdErr := pim.StatCommand(process, &database)
  if cmdErr != nil {
    t.Errorf("Error in info command: %v", cmdErr)
  }
}
