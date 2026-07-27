package main

import (
	"flag"
	"fmt"
	"os"
)

type DFPMArgs struct {
	Command string
	Args    []string

	// Flags

	PackDir string // install, list, search, xc8

	Download bool // search
}

var args DFPMArgs
var scmd *flag.FlagSet

const (
	sh = " (shorthand)"

	defaultPackDir = "~/.mchp_packs"
	usagePackDir   = "pack installation directory"

	defaultDownload = false
	usageDownload   = "download pack(s) automatically"
)

func InitCLI() {

	flag.Usage = usage

	flag.BoolFunc("version", "", version)

	flag.StringVar(&args.PackDir, "pack-dir", defaultPackDir, usagePackDir)
	flag.StringVar(&args.PackDir, "p", defaultPackDir, usagePackDir+sh)

	scmd = flag.NewFlagSet("search", flag.ExitOnError)

	scmd.BoolVar(&args.Download, "download", defaultDownload, usageDownload)
	scmd.BoolVar(&args.Download, "d", defaultDownload, usageDownload+sh)
}

func ParseCLI() DFPMArgs {
	flag.Parse()

	args.PackDir = ExpandPath(args.PackDir)

	if flag.Arg(0) == "" {
		usageErr()
	}

	switch flag.Arg(0) {
	case "search":
		scmd.Parse(flag.Args()[1:])
		args.Args = scmd.Args()
	default:
		args.Args = flag.Args()[1:]
	}
	args.Command = flag.Arg(0)
	return args
}

func usage() {
	fmt.Println(`dfpm - Device Family Pack Manager

Usage:
  dfpm [-p <path>] install <pack>
  dfpm [-p <path>] list
  dfpm [-p <path>] search [-d] <string>
  dfpm [-p <path>] xc8 [-mcpu=<cpu_id>] [<argument>...]

Options:
  -p <path>, --pack-dir <path>  specify pack directory to use [default: ~/.mchp_packs]
  -d, --download                automatically download packs
  -h, --help                    print help and exit
  --version                     print version and exit`)
}

func usageErr() {
	usage()
	os.Exit(1)
}

func version(flag string) error {
	fmt.Printf("dfpm - Device Family Pack Manager v%s\nCopyright (c) %s Rex McKinnon\n", Dfpm_Version, Dfpm_Year)
	os.Exit(0)
	return nil
}
