package lib

import (
  "context"
  "os/exec"
  "strings"
  "syscall"
  "os"
  // "errors"
  "time"
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
  // var exitErr *exec.ExitError
  app, args := formatCommand(run.Command)
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  cmd := exec.CommandContext(ctx, app, args...)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
  cmd.Dir = run.Directory
  cmd.Env = os.Environ()
  defer cancel()
  if runErr := cmd.Run(); runErr != nil {
    syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
    return "terminator!", 0
  }
  return "not terminated", 1

  // if showOutput {
  //   Info.Printf("%s", string(runOutput))
  // }
  // if errors.As(runErr, &exitErr) {
  //   return string(runOutput), exitErr.ExitCode()
  // } else if runErr != nil {
  //   return string(runOutput), -1
  // }
  // return string(runOutput), 0
}

