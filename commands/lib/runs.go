package lib

import (
  "context"
  "os/exec"
  "strings"
  "syscall"
  "errors"
  "time"
  "os"
)

func generateLogName(length int) string {
  id := make([]byte, length)
  for i := range id {
    id[i] = IDCHARSET[SeededRand.Intn(len(IDCHARSET))]
  }
  return CONFIGDIR + "/" + string(id) + ".log"
}

func getCommandLogs(filename string) string {
  content, fileErr := os.ReadFile(filename)
  cmd := exec.Command("rm", filename)
  defer cmd.Run() //nolint:errcheck
  if fileErr != nil {
    Error.Printf("Can't read from temp log file %s", filename)
    return ""
  }
  return string(content)
}

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

func ExecuteTimedRun(run Run, showOutput bool, duration int) (string, int) {
  logName := generateLogName(5)
  log, _ := os.Create(logName)
	defer log.Close()
  app, args := formatCommand(run.Command)
  ctx, _ := context.WithTimeout(context.Background(), time.Duration(duration) * time.Second) //nolint:govet
  cmd := exec.CommandContext(ctx, app, args...)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
  cmd.Dir = run.Directory
  cmd.Env = os.Environ()
	cmd.Stdout = log
  if runErr := cmd.Run(); runErr != nil {
    syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) //nolint:errcheck
    return getCommandLogs(logName), 0
  }
  return "not terminated", 1
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

