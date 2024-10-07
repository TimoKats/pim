package lib

import (
  "math/rand"
  "errors"
  "time"
)

var (
  SeededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func generateId(length int) string {
  id := make([]byte, length)
  for i := range id {
    id[i] = IDCHARSET[SeededRand.Intn(len(IDCHARSET))]
  }
  return string(id)
}

func StoreRun(run Run, output string, exitCode int) Log {
  runTimeStamp := time.Now()
  return Log {
    Id: generateId(5),
    RunCommand: run,
    Output: output,
    ExitCode: exitCode,
    Timestamp: runTimeStamp,
  }
}

// maybe a seperate module for this? idk.

func ViewLog (database *Database, logId string) error {
  for _, log := range database.Logs {
    if log.Id == logId {
      timeString := log.Timestamp.Format(time.RFC822Z)
      Info.Printf("exit code: %d", log.ExitCode)
      Info.Printf("command: %s", log.RunCommand.Command)
      Info.Printf("timestamp: %s", timeString)
      Info.Printf("directory: %s", log.RunCommand.Directory)
      Info.Printf("schedule: %s", log.RunCommand.Schedule)
      Info.Println("command output:\n---")
      Info.Println(log.Output)
      return nil
    }
  }
  return errors.New("Log id not found in data.")
}

func ViewLogs (database *Database) error {
  for _, log := range database.Logs {
    timeString := log.Timestamp.Format(time.RFC822Z)
    id := ResponsiveWhitespace(log.Id)
    name := ResponsiveWhitespace(log.RunCommand.Name)
    Info.Printf("%d | %s | %s | %s", log.ExitCode, name, timeString, id)
  }
  return nil
}


