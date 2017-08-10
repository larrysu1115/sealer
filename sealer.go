package main

import (
	"fmt"
	"os"
)

func main() {
	flags, err := LoadFlags(os.Args[1:])
	if err != nil {
		fmt.Printf("err initial flags: %v", err)
		return
	}
	SetupLogger(flags.LogFile, "INFO", `%{time:060102 15:04:05} - %{level:.4s} %{shortfunc} %{message}`)
	// Lg.Infof("flags: %v", flags)

	switch flags.Command {
	case "status":
		cmdStatus(flags)
	case "do":
		err = cmdDoArchive(flags)
	case "undo":
		err = cmdDoUnarchive(flags)
	default:
		Lg.Warningf("Unknown Command:[%s]! Exit.", flags.Command)
	}

	if err != nil {
		Lg.Errorf("Failed! %v", err)
	}
}
