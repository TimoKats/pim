// No tedious unit testing. Just run some of the main commands and see if anything breaks.
// Note, I want to add start here but I don't have a file system in GH actions for the
// lock files? So that fix is still coming.

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
    lib.Run{
      Name: "sleepy-command",
      Schedule: "@start+10",
      Command: "echo bye world",
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

