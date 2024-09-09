// Contains commands that run every startup automatically. In short, the data file and
// yaml files are loaded (sometimes created) and parsed into objects through the
// SetupYamlFiles function. Next, certain attributes are cleaned and checked by helper
// functions. Also, there the CheckStartupErrors is run every startup to check if all
// const values that are needed can be set.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "strings"
  "errors"
)

func formatProcessName(processName string) string {
  processName = strings.ToLower(processName)
  processName = strings.Replace(processName, " ", "-", -1)
  processName = strings.Replace(processName, "_", "-", -1)
  return processName
}

func formatProcess(process *lib.Process) {
  var processName string
  for index, _ := range process.Runs {
    processName = process.Runs[index].Name
    process.Runs[index].Name = formatProcessName(processName)
  }
}

func SetupYamlFiles() (lib.Process, lib.Database, error) {
  database, readDataErr := lib.ReadDataYaml()
  process, readProcessErr := lib.ReadProcessYaml()
  formatProcess(&process)
  return process, database, errors.Join(readDataErr, readProcessErr)
}

func CheckStartupErrors() error {
  if err := errors.Join(lib.CONFIGERR, lib.LOGERR); err != nil {
    return err
  }
  return nil
}

