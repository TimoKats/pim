package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
)

func SetupYamlFiles() (lib.Process, lib.Database, error) {
  database, readDataErr := lib.ReadDataYaml(lib.DATAPATH) // NOTE: Pass them as variable?
  process, readProcessErr := lib.ReadProcessYaml(lib.PROCESSPATH)
  return process, database, errors.Join(readDataErr, readProcessErr)
}
