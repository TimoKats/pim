package lib

import (
  "time"
)

func StoreRun(run Run, output string, exitCode int) Log {
  runTimeStamp := time.Now()
  return Log {
    RunCommand: run,
    Output: output,
    ExitCode: exitCode,
    Timestamp: runTimeStamp,
  }
}

