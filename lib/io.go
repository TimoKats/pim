package lib

import (
  "gopkg.in/yaml.v2"
  "errors"
  "os"
)

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
  return process, fileErr
}

func WriteDataYaml(filename string, database Database) error {
  yamlData, yamlErr := yaml.Marshal(&database)
  writeErr := os.WriteFile(filename, yamlData, 0644)
  if err := errors.Join(yamlErr, writeErr); err != nil {
    return err
  }
  return nil
}

