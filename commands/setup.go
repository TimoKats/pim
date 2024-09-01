package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "strings"
  "errors"
)

func formatProcessName(processName string) string { // NOTE: move to lib?
  processName = strings.ToLower(processName)
  processName = strings.Replace(processName, " ", "-", -1)
  processName = strings.Replace(processName, "_", "-", -1)
  return processName
}

func formatProcess(process *lib.Process) { // NOTE: move to lib?
  var processName string
  for index, _ := range process.Runs {
    processName = process.Runs[index].Name
    process.Runs[index].Name = formatProcessName(processName)
  }
}

func SetupYamlFiles() (lib.Process, lib.Database, error) {
  database, readDataErr := lib.ReadDataYaml(lib.DATAPATH) // NOTE: Pass them as variable?
  process, readProcessErr := lib.ReadProcessYaml(lib.PROCESSPATH)
  formatProcess(&process)
  return process, database, errors.Join(readDataErr, readProcessErr)
}

