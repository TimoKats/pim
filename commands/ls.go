// Lists the processes in the current process.yaml.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"
)

func ListCommand(process lib.Process) error {
  for _, run := range process.Runs {
    name := lib.ResponsiveWhitespace(run.Name)
    cmd := lib.ResponsiveWhitespace(run.Command)
    schedule := lib.ResponsiveWhitespace(run.Schedule)
    lib.Info.Printf("%s | %s | %s | %d ", name, schedule, cmd, run.Duration)
  }
  return nil
}

