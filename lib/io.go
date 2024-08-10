package lib

import (
  "gopkg.in/yaml.v2"
  "io/ioutil"
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

func WriteDataYaml(filename string, database Database) {
  yamlData, _ := yaml.Marshal(&database)
  ioutil.WriteFile(filename, yamlData, 0644)
} // NOTE: add error handling here!

