package main

import (
  lib "github.com/TimoKats/pim/commands/lib"
  pim "github.com/TimoKats/pim/commands"

  "testing"
)

var process lib.Process
var database lib.Database

func TestSetup(t *testing.T) {
  processTemp, databaseTemp, setupErr := pim.SetupYamlFiles()
  lib.Info.Printf("HELLO TIMO! %s", lib.CONFIGDIR)
  if setupErr != nil {
    t.Errorf("Error in setup: %v", setupErr)
  }
  process = processTemp
  database = databaseTemp
}

func TestLs(t *testing.T) {
  cmdErr := pim.ListCommand(process)
  if cmdErr != nil {
    t.Errorf("Error in ls command: %v", cmdErr)
  }
}

func TestRun(t *testing.T) {
  command := []string{"one", "sleep", "python-version"}
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
