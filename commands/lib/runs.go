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
  return LOGDIR + "/" + string(id) + ".log"
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

func executeTimedRun(run Run, showOutput bool, duration int) (string, int) {
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

func executeRun(run Run, showOutput bool) (string, int) {
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

func ExecuteCommand(command string) (string, int) {
  var exitErr *exec.ExitError
  app, args := formatCommand(command)
  cmd := exec.Command(app, args...)
  runOutput, runErr := cmd.CombinedOutput()
  if errors.As(runErr, &exitErr) {
    return string(runOutput), exitErr.ExitCode()
  } else if runErr != nil {
    return string(runOutput), -1
  }
  return string(runOutput), 0
}

func RunAndStore(run Run, database *Database, process Process, showOutput bool) {
  var output string
  var status int
  if run.Duration != 0 {
    output, status = executeTimedRun(run, showOutput, run.Duration)
  } else {
    output, status = executeRun(run, showOutput)
  }
  storedLog := StoreRun(run, output, status)
  database.Logs = append(database.Logs, storedLog)
  if (!process.OnlyStoreErrors) || (process.OnlyStoreErrors && status != 0) {
    writeErr := WriteDataYaml(DATAPATH, *database)
    if writeErr != nil {
      Warn.Printf("process %s is not stored in database.", run.Name)
    }
  }
}
