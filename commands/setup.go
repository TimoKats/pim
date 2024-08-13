package commands

import (
  lib "pim/commands/lib"

  "errors"
)

func SetupYamlFiles() (lib.Process, lib.Database, error) {
  database, readDataErr := lib.ReadDataYaml("tests/data.yaml")
  process, readProcessErr := lib.ReadProcessYaml("tests/test.yaml")
  return process, database, errors.Join(readDataErr, readProcessErr)
}
