// Functions needed to store items in the database. Database itself is loaded/written in
// the io module.

package lib

import (
  "math/rand"
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

