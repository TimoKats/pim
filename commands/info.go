// Triggered by the pim info command. Shows information related to this installation of
// pim and its creator.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"
)

func Info() error {
  lib.Info.Println("Description : Pim is a process orchestrator.")
  lib.Info.Printf("Version     : %s", lib.VERSION)
  lib.Info.Printf("Author      : %s", lib.AUTHOR)
  lib.Info.Printf("License     : %s", lib.LICENSE)
  lib.Info.Println("More info   : https://github.com/TimoKats/pim")
  return nil
}

