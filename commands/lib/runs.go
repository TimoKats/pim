package lib

import (
  "os/exec"
  "strings"
  "errors"
)

func formatCommand(command string) (string, []string)  {
  var app string
  var args []string
  for index, text := range strings.Fields(command) {
    if index == 0 {
      app = text
    } else {
      args = append(args, text)
    }
  }
  return app, args
}

func ExecuteRun(run Run, showOutput bool) (string, int) {
  var exitErr *exec.ExitError
  app, args := formatCommand(run.Command)
  cmd := exec.Command(app, args...)
  cmd.Dir = run.Directory
  runOutput, runErr := cmd.CombinedOutput()
  if showOutput {
    Info.Printf("%s", string(runOutput))
  }
  if errors.As(runErr, &exitErr) {
    return string(runOutput), exitErr.ExitCode()
  } else if runErr != nil {
    return string(runOutput), -1
  }
  return string(runOutput), 0
}

