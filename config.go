package main

import (
	"fmt"

	ini "gopkg.in/ini.v1"

	flags "github.com/jessevdk/go-flags"
)

// CommandFlags is flags from command line.
type CommandFlags struct {
	Command       string `long:"cmd"           description:"Command to run"   required:"true"`
	DirSrc        string `long:"src"           description:"source dir"       required:"true"`
	DirTo         string `long:"to"            description:"destination dir"  required:"true"`
	Prefix        string `long:"pre"           description:"archive prefix"   `
	FileConfig    string `long:"cfg"           description:"config file"      default:"./sealer.ini"`
	FilePattern   string `long:"file"          description:"regex filename"   default:""`
	PackNumber    int64  `long:"packNumber"    description:"number of files in one tgz" default:"3000"`
	RetainSeconds int64  `long:"retainSeconds" description:"retain seconds"   default:"172800"`
	LogFile       string `long:"logfile"       description:"logfile path"     default:""`
}

func (o CommandFlags) String() string {
	return fmt.Sprintf("<%s>, DirSrc=[%s], DirTo=[%s], FilePattern=[%s] Pack=%d, RetainSeconds=%d, Logfile=[%s]",
		o.Command, o.DirSrc, o.DirTo, o.FilePattern, o.PackNumber, o.RetainSeconds, o.LogFile)
}

// LoadFlags reads config.ini file.
func LoadFlags(args []string) (f CommandFlags, err error) {
	_, err = flags.NewParser(&f, flags.Default|flags.IgnoreUnknown).ParseArgs(args)
	if err != nil {
		fmt.Printf("Failed parsing command line flags. %v", err)
		return
	}
	if f.FileConfig == "" {
		return
	}

	cfg, err := ini.Load(f.FileConfig)
	if err != nil {
		fmt.Printf("Cannot load config file: [%s], err: %v", f.FileConfig, err)
		return
	}
	cfgSec := cfg.Section("sealer")
	if cfgSec.HasKey("file_pattern") && f.FilePattern == "" {
		f.FilePattern = cfgSec.Key("file_pattern").String()
	}
	if cfgSec.HasKey("retain_seconds") {
		f.RetainSeconds = cfgSec.Key("retain_seconds").MustInt64(259200)
	}
	if cfgSec.HasKey("pack_number") {
		f.PackNumber = cfgSec.Key("pack_number").MustInt64(1000)
	}
	if cfgSec.HasKey("log_file") {
		f.LogFile = cfgSec.Key("log_file").String()
	}
	return
}
