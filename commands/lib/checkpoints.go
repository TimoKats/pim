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
        Next: job.NextRun(),
        Name: runs[index].Name,
        Catchup: runs[index].Catchup,
      })
  }
  return Checkpoint{Updated: time.Now(), Runs: checkpoints}
}

func WriteCheckpoint(runs []Run, schedule *gocron.Scheduler) error {
  checkpoint := CreateCheckpoint(runs, schedule)
  yamlData, yamlErr := yaml.Marshal(&checkpoint)
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

