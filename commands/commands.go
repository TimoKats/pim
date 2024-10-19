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
  lockErr := os.Remove(lib.LOCKPATH)
  process, processErr := os.FindProcess(intPid)
  if err := errors.Join(lockErr, processErr); err != nil {
    return processErr
  }
  lib.Warn.Printf("Stopping pim process %d", intPid)
  killErr := process.Kill()
  return killErr
}

func ListCommand(process lib.Process, database *lib.Database) error {
  dummySchedule := lib.CreateDummySchedule(process, database)
  lib.Info.Println(lib.ViewListHeader())
  for _, run := range process.Runs {
    nextRun, _ := lib.ViewNextRun(dummySchedule, run)
    name := lib.ResponsiveWhitespace(run.Name)
    cmd := lib.ResponsiveWhitespace(run.Command)
    schedule := lib.ResponsiveWhitespace(run.Schedule)
    lib.Info.Printf("%s | %s | %s | %s", name, schedule, cmd, nextRun)
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

func StatusCommand() error {
  pid, lockErr := lib.ReadLockFile()
  processCount := lib.CountPimProcesses()
  if processCount > 0 {
    lib.Info.Printf("Pim is currently running at: %d", pid)
    return lockErr
  } else {
    lib.Info.Println("No pim process running.")
    return nil
  }
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
      lib.RunJobMapping[cronJob] = run
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

