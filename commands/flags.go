// Obviously not a command, but if user submits a flag then this module is used to parse
// it. This module is there mainly to keep main package clean and adhere to my design.

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
)

const asciiLogo = `
     _         | Description: Pim is a task orchestration tool.
 ___|_|_____   | Version: %s
| . | |     |  | Author: %s
|  _|_|_|_|_|  | License: %s
|_|            | Source code: github.com/TimoKats/pim

`

func FlagCommand(flag string) error {
  switch flag {
    case "version":
      lib.Info.Println(lib.VERSION)
      return nil
    case "help":
      lib.Info.Println(lib.ABSTRACT)
      lib.Info.Println(lib.HELPSTRING)
      return nil
    case "info":
      lib.Info.Printf(asciiLogo, lib.VERSION, lib.AUTHOR, lib.LICENSE)
      return nil
    case "license":
      lib.Info.Println(lib.LICENSE)
      return nil
    default:
      return errors.New("flag '" + flag + "' not found.")
  }
}

