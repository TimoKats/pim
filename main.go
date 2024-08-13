package main

import (
  pim "pim/lib"
  "errors"
  "time"
  "os"
  "github.com/robfig/cron"
)

func heatbeat() {
  for {
    time.Sleep(time.Second)
  }
}

func runRun (run pim.Run, database *pim.Database) {
  pim.Info.Printf("Now running '%s'", run.Name)
  output, status := pim.ExecuteRun(run)
  storedLog := pim.StoreRun(run, output, status)
  database.Logs = append(database.Logs, storedLog)
  pim.WriteDataYaml("tests/data.yaml", *database)
}


func runSelected (selectedRun string, process pim.Process, database *pim.Database) error {
  for _, run := range process.Runs {
    if run.Name == selectedRun {
      output, status := pim.ExecuteRun(run)
      storedLog := pim.StoreRun(run, output, status)
      database.Logs = append(database.Logs, storedLog)
      pim.WriteDataYaml("tests/data.yaml", *database)
      return nil
    }
  }
  return errors.New("Name of selected run not in process yaml.")
}

func runSchedule(process pim.Process, database *pim.Database) {
  cronSchedule := cron.New()
  for _, run := range process.Runs {
    run := run
    cronErr := cronSchedule.AddFunc(run.Schedule, func() { runRun(run, database) })
    if cronErr != nil {
      pim.Error.Printf("Error in '%s'. Check Yaml.", run.Name)
    }
  }
  cronSchedule.Start()
}

func parseCommand(command []string, process pim.Process, database *pim.Database) error  {
  switch command[1] {
    case "run":
      if len(command) < 3 {
        return errors.New("No command name given. pim run <<name>>.")
      }
      return runSelected(command[2], process, database)
    }
  return nil
}

func main() {
  if len(os.Args) < 2 {
    pim.Error.Println("Not enough arguments. pim <<run, start, ls>>.")
    return
  }
  database, readDataErr := pim.ReadDataYaml("tests/data.yaml")
  process, readProcessErr := pim.ReadProcessYaml("tests/test.yaml")
  if err := errors.Join(readDataErr, readProcessErr); err != nil {
    return
  }
  parseErr := parseCommand(os.Args, process, &database)
  pim.Info.Println(parseErr)
  // runSchedule(process, &database)
  // heatbeat()
}
