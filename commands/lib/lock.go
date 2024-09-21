package lib

import (
  "strconv"
  "os/exec"
  "strings"
  "errors"
  "os"
)

func RemoveLockFile() error {
  cmd := exec.Command("rm", LOCKPATH)
  return cmd.Run() //nolint:errcheck
}

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

func RemoveDanglingLock() {
  processCount := 0
  test, _ := ExecuteCommand("ps -u")
  for _, line := range strings.Split(test, "\n") {
    if strings.Contains(line, "pim start") {
      processCount += 1
    }
  }
  if processCount < 2 && LockExists() {
    removeErr := RemoveLockFile()
    if removeErr != nil {
      Error.Println("Failed removing dangling lock file.")
    }
  }
}

