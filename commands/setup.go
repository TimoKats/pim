package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
)

func SetupYamlFiles() (lib.Process, lib.Database, error) {
  database, readDataErr := lib.ReadDataYaml("/home/timokats/.pim/data.yaml")
  process, readProcessErr := lib.ReadProcessYaml("/home/timokats/.pim/process.yaml")
  return process, database, errors.Join(readDataErr, readProcessErr)
}
