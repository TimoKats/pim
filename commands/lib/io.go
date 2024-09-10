package lib

import (
  "gopkg.in/yaml.v2"
  "errors"
  "os"
)

func ReadDataYaml() (Database, error) {
  var database Database
  yamlFile, fileErr := os.ReadFile(DATAPATH)
  if fileErr == nil {
    yamlErr := yaml.Unmarshal(yamlFile, &database)
    if yamlErr == nil {
      return database, nil
    }
    return database, yamlErr
  } else if errors.Is(fileErr, os.ErrNotExist) {
    Warn.Printf("Creating new data file at %s", PROCESSPATH)
    _, createErr := os.Create(PROCESSPATH)
    return database, createErr
  }
  return database, fileErr
}

func ReadProcessYaml() (Process, error) {
  var process Process
  yamlFile, fileErr := os.ReadFile(PROCESSPATH)
  if fileErr == nil {
    yamlErr := yaml.Unmarshal(yamlFile, &process)
    if yamlErr == nil {
      return process, nil
    }
    return process, yamlErr
  }
  Warn.Printf("No database of runs found. Please create one at: %s", PROCESSPATH)
  return process, nil
}

func WriteDataYaml(filename string, database Database) error {
  yamlData, yamlErr := yaml.Marshal(&database)
  writeErr := os.WriteFile(filename, yamlData, 0644)
  if err := errors.Join(yamlErr, writeErr); err != nil {
    return err
  }
  return nil
}

func TrimDatabase(database *Database, threshold int) {
  if len(database.Logs) > threshold && threshold != 0 {
    database.Logs = database.Logs[len(database.Logs) - threshold:]
  }
  writeErr := WriteDataYaml(DATAPATH, *database)
  if writeErr != nil {
    Warn.Println("Wasn't able to trim database. Continuing operations...")
  }
}

