package lib

import (
  "errors"
  "time"
  "os"

  "github.com/go-co-op/gocron"
  "gopkg.in/yaml.v2"
)

func CreateCheckpoint(runs []Run, schedule *gocron.Scheduler) Checkpoint {
  var checkpoints []RunCheckpoint
  for index, job := range schedule.Jobs() {
    checkpoints = append(
      checkpoints,
      RunCheckpoint{
        NextRun: job.NextRun(),
        Name: runs[index].Name,
      })
  }
  return Checkpoint{Test: "hello", RunCheckpoints: checkpoints}
}

func WriteCheckpoint(runs []Run, schedule *gocron.Scheduler) error {
  checkpoint := CreateCheckpoint(runs, schedule) // NOTE: I AM HERE
  yamlData, yamlErr := yaml.Marshal(&checkpoint)
  Warn.Printf("> %v", yamlData)
  writeErr := os.WriteFile(CHECKPOINTPATH, yamlData, 0644)
  if err := errors.Join(yamlErr, writeErr); err != nil {
    return err
  }
  return nil
}

func RemoveCheckpoint() error {
  removeErr := os.Remove(CHECKPOINTPATH)
  return removeErr
}

func ReadCheckpoint() (time.Time, error) {
  byteCheckpoint, readErr := os.ReadFile(CHECKPOINTPATH)
  if readErr != nil {
    return time.Now(), readErr
  }
  strCheckpoint := string(byteCheckpoint)
  timeCheckpoint, convertErr := time.Parse(time.RFC850, strCheckpoint)
  return timeCheckpoint, convertErr
}

