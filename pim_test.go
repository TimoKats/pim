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
      Duration: 5,
    },
  },
}


func TestLs(t *testing.T) {
  cmdErr := pim.ListCommand(process, &database)
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

