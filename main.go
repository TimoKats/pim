package main

import (
  pim "pim/lib"
)

func main() {
  process, readYamlErr := pim.ReadProcessYaml("tests/test.yaml")
  var database pim.Database
  if readYamlErr == nil {
    for _, run := range process.Runs {
      output, status := pim.ExecuteRun(run)
      storedLog := pim.StoreRun(run, output, status)
      database.Logs = append(database.Logs, storedLog)
    }
    pim.WriteDataYaml("tests/data.yaml", database)
  }
}
