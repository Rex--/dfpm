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
	fmt.Println(`usage: dfpm [-p <pack_dir>] [install | list | search] [arg]`)
	flag.PrintDefaults()
	fmt.Println(`
  dfpm install <pack>      install the given device pack
  dfpm list                list installed device packs
  dfpm search <string>     search for <string> in pack descriptions

  dfpm xc8 [-mcpu=<cpuid>] [args ...]
                           search for the given <cpuid> in pack descriptions
                           if found, launch 'xc8-cc' with the -mdfp option`)
}

func usageErr() {
	usage()
	os.Exit(1)
}
