// Flags can be called from main package using a - or -- prefix. Currently, the Flags
// only output constant information (like license, version, etc) to standard output. The
// only reason they are in their own submodule is because it keeps other modules cleaner
// :)

package commands

import (
  lib "github.com/TimoKats/pim/commands/lib"

  "errors"
)

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
      lib.Info.Printf(lib.ASCIILOGO, lib.VERSION, lib.AUTHOR, lib.LICENSE)
      return nil
    case "license":
      lib.Info.Println(lib.LICENSE)
      return nil
    default:
      return errors.New("flag '" + flag + "' not found.")
  }
}

