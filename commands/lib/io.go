package lib

import (
  "gopkg.in/yaml.v2"
  "errors"
  "os"
)

func ReadDataYaml(filename string) (Database, error) {
  var database Database
  yamlFile, fileErr := os.ReadFile(filename)
  if fileErr == nil {
    yamlErr := yaml.Unmarshal(yamlFile, &database)
    if yamlErr == nil {
      return database, nil
    }
    return database, yamlErr
  } else if errors.Is(fileErr, os.ErrNotExist) {
    Warn.Printf("Creating new data file at %s", filename)
    os.Create(filename)
    return database, nil
  }
  return database, fileErr
}

func ReadProcessYaml(filename string) (Process, error) {
  var process Process
  yamlFile, fileErr := os.ReadFile(filename)
  if fileErr == nil {
    yamlErr := yaml.Unmarshal(yamlFile, &process)
    if yamlErr == nil {
      return process, nil
    }
    return process, yamlErr
  }
  Warn.Printf("No database of runs found. Will add new file at: %s", filename)
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

