// contains all the functions that can be invoked from the main package.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "strings"
  "errors"
  "time"
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

func ListCommand(process lib.Process) error {
  for _, run := range process.Runs {
    name := lib.ResponsiveWhitespace(run.Name)
    cmd := lib.ResponsiveWhitespace(run.Command)
    schedule := lib.ResponsiveWhitespace(run.Schedule)
    lib.Info.Printf("%s | %s | %s | %d ", name, schedule, cmd, run.Duration)
  }
  return nil
}

func RunCommand(command []string, process lib.Process, database *lib.Database) error {
  var selectedRun string
  if len(command) < 3 {
    return errors.New("No command name given. pim run <<name>>.")
  } else {
    selectedRun = command[2]
  }
  for _, run := range process.Runs {
    if run.Name == selectedRun {
      lib.RunAndStore(run, database, process, true)
      return nil
    }
  }
  return errors.New("'" + command[2] + "' not in process yaml.")
}

func LogCommand(command []string, database *lib.Database) error {
  if len(command) < 3 {
    return lib.ViewLogs(database)
  }
  return lib.ViewLog(database, command[2])
}

func StartCommand(process lib.Process, database *lib.Database) error {
  if setupErr := SetupStart(); setupErr != nil {
    return setupErr
  }
  for _, run := range process.Runs {
    run := run
    cronJob, cronErr := lib.SelectCron(run, process, database)
    if cronErr != nil {
      lib.Error.Printf("Error in '%s'. %v.", run.Name, cronErr)
    } else if cronJob != nil {
      cronJob.Tag(run.Name)
    }
  }
  lib.Schedule.StartAsync()
  lib.Catchup()
  lib.Heartbeat(process, database)
  return nil
}

func CleanCommand(database *lib.Database) error {
  database.Logs = nil
  writeErr := lib.WriteDataYaml(lib.DATAPATH, *database)
  if writeErr != nil {
    return writeErr
  }
  os.RemoveAll(lib.LOGDIR)
  makeErr := os.Mkdir(lib.LOGDIR, 0755)
  return makeErr
}

func TestCommand(process lib.Process, database *lib.Database) error {
  var nextRun string
  var runsCatchup bool
  lib.Info.Printf("%s | %s | %s", lib.ResponsiveWhitespace("Name"), lib.ResponsiveWhitespace("Next scheduled run"), "Runs on startup")
  for _, run := range process.Runs {
    run := run
    cronJob, cronErr := lib.DummyCron(run, process, database)
    if cronJob != nil && cronErr == nil {
      cronJob.Tag(run.Name)
    }
  }
  lib.Schedule.StartAsync()
  for _, run := range process.Runs {
    cronJob, cronErr := lib.Schedule.FindJobsByTag(run.Name)
    if cronErr == nil && cronJob != nil {
      nextRun = lib.ResponsiveWhitespace(cronJob[0].NextRun().Format(time.RFC1123))
      runsCatchup = lib.RunsCatchup(run.Name)
    } else {
      nextRun = lib.ResponsiveWhitespace("No cron schedule.")
      runsCatchup = strings.HasPrefix(run.Schedule, "@start")
    }
    runName := lib.ResponsiveWhitespace(run.Name)
    lib.Info.Printf("%s | %s | %v", runName, nextRun, runsCatchup)
  }
  return nil
}
