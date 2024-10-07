// contains all the functions that can be invoked from the main package.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
  "os"
)

// Used to end runs started by pim start (assuming you run it in the background). It
// checks if the lock file exists and returns the error based on the PID in the lock
// file. Note, it doesn't hurt to occasionally run ps -aux | grep pim to double check
// if anything else is running.
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

// Lists the processes in the current process.yaml.
func ListCommand(process lib.Process) error {
  for _, run := range process.Runs {
    name := lib.ResponsiveWhitespace(run.Name)
    cmd := lib.ResponsiveWhitespace(run.Command)
    schedule := lib.ResponsiveWhitespace(run.Schedule)
    lib.Info.Printf("%s | %s | %s | %d ", name, schedule, cmd, run.Duration)
  }
  return nil
}

// Run a command from process.yaml based on its name (which is supplied as an argument)
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

// Shows logs of previous runs. If pim log is run, then a summary/table is shown of
// all previous runs (also shows run-id). if pim log <run-id> is run, then a more
// elaborate overview is shown with ViewLog.
func LogCommand(command []string, database *lib.Database) error {
  if len(command) < 3 {
    return lib.ViewLogs(database)
  }
  return lib.ViewLog(database, command[2])
}

// This function uses an external cron module that runs different jobs concurrently.
// Update: I create a lock file that contains the current pid (of the start command)
// that is meant to prevent multiple start processes running at the same time and kill
// the process without calling ps -aux > kill ...
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

// Contains functions related to cleaning/trimming the files that pim creates. First, the
// the function TrimDatabase is run every heartbeat/run and it's behavior is determined
// by the set_max_logs setting in your process.yaml. Next, the CleanDatabase function
// deletes all the redundant log files and previous runs in data.yaml.
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


