// Commands that run every startup automatically. In short, the data file and yaml files
// are loaded (or created; if they don't exist yet) and parsed into objects through the
// SetupYamlFiles function. Next, certain attributes are cleaned and checked by helper
// functions. Finally, the CheckStartupErrors is called every startup to check if all
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
  for i := 0; i < len(process.Runs); i++ {
    processName = process.Runs[i].Name
    process.Runs[i].Name = formatProcessName(processName)
  }
}

func SetupStart() error {
  lib.InitFileLogging()
  lib.RemoveDanglingLock()
  if lib.LockExists() {
    return errors.New("Pim is already running! Run <<pim stop>> or check lockfile/ps.")
  }
  lockErr := lib.InitLockFile()
  if lockErr != nil {
    return lockErr
  }
  return nil
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

