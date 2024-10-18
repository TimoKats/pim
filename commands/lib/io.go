// All reading/writing operations to the filesystem are in this submodule. The categories
// of functions are: data.yaml, process.yaml, checkpoints, lockfiles.

package lib

import (
  "github.com/go-co-op/gocron"
  "gopkg.in/yaml.v2"

  "strconv"
  "strings"
  "errors"
  "time"
  "os"
)

// data.yaml

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
    Warn.Printf("Creating new data file at %s", DATAPATH)
    _, createErr := os.Create(DATAPATH)
    return database, createErr
  }
  return database, fileErr
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

// process.yaml

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

// lockfiles

func InitLockFile() error {
  currentPid := strconv.Itoa(os.Getpid())
  lockErr := os.WriteFile(LOCKPATH, []byte(currentPid), 0644)
  return lockErr
}

func LockExists() bool {
  if _, lockErr := os.Stat(LOCKPATH); errors.Is(lockErr, os.ErrNotExist) {
    return false
  }
  return true
}

func ReadLockFile() (int, error) {
  bytePid, lockErr := os.ReadFile(LOCKPATH)
  if lockErr != nil {
    return 0, lockErr
  }
  strPid := string(bytePid)
  intPid, convErr := strconv.Atoi(strPid)
  if convErr != nil {
    return 0, convErr
  }
  return intPid, nil
}

func CountPimProcesses() int {
  processCount := 0
  test, _ := ExecuteCommand("ps -u")
  for _, line := range strings.Split(test, "\n") {
    if strings.Contains(line, "pim start") {
      processCount += 1
    }
  }
  return processCount
}

func RemoveDanglingLock() {
  processCount := CountPimProcesses()
  if processCount < 2 && LockExists() {
    removeErr := os.Remove(LOCKPATH)
    if removeErr != nil {
      Error.Println(removeErr)
      Error.Println("Failed removing dangling lock file.")
    }
  }
}

// checkpoints

func CreateCheckpoint(schedule *gocron.Scheduler, mapping map[*gocron.Job]Run) Checkpoint {
  var checkpoints []RunCheckpoint
  for _, job := range schedule.Jobs() {
    checkpoints = append(
      checkpoints,
      RunCheckpoint{
        Next: job.NextRun(),
        Name: mapping[job].Name,
        Catchup: mapping[job].Catchup,
      })
  }
  return Checkpoint{Updated: time.Now(), Runs: checkpoints}
}

func WriteCheckpoint(schedule *gocron.Scheduler, mapping map[*gocron.Job]Run) error {
  checkpoint := CreateCheckpoint(schedule, mapping)
  yamlData, yamlErr := yaml.Marshal(&checkpoint)
  writeErr := os.WriteFile(CHECKPOINTPATH, yamlData, 0644)
  if err := errors.Join(yamlErr, writeErr); err != nil {
    return err
  }
  return nil
}

func ReadCheckpoint() (Checkpoint, error) {
  var checkpoint Checkpoint
  yamlFile, fileErr := os.ReadFile(CHECKPOINTPATH)
  if fileErr == nil {
    yamlErr := yaml.Unmarshal(yamlFile, &checkpoint)
    if yamlErr == nil {
      return checkpoint, nil
    }
    return checkpoint, yamlErr
  }
  return checkpoint, errors.New("No checkpoint found")
}
