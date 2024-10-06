// Contains the stop command. Used to end runs started by pim start (assuming you run it
// in the background). It checks if the lock file exists and returns the error based on
// the PID in the lock file. Note, it doesn't hurt to occasionally run ps -aux | grep pim
// to double check if anything else is running.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
  "os"
)

func StopCommand() error {
  if !lib.LockExists() {
    return errors.New("No process to end.")
  }
  intPid, readErr := lib.ReadLockFile()
  if readErr != nil {
    return readErr
  }
  lockErr := os.Remove(lib.CHECKPOINTPATH)
  process, processErr := os.FindProcess(intPid)
  if err := errors.Join(lockErr, processErr); err != nil {
    return processErr
  }
  lib.Warn.Printf("Stopping pim process %d", intPid)
  killErr := process.Kill()
  return killErr
}

