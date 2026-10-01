package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/integrii/flaggy"
)

var Dfpm_Version = "1.1.0"

// Subcommand vars
var (
	cmdInstall *flaggy.Subcommand
	cmdRemove  *flaggy.Subcommand
	cmdList    *flaggy.Subcommand
	cmdSearch  *flaggy.Subcommand
	cmdUpdate  *flaggy.Subcommand
	cmdXc8     *flaggy.Subcommand
)

// Flag vars
var (
	packDir       string = "~/.mchp_packs/Microchip"
	downloadPacks bool   = false
	yesDownload   bool   = false
	yesRemove     bool   = false
	localSearch   bool   = false
	listAll       bool   = false
	showUpdates   bool   = false
	keepPacks     bool   = false
)

// Positional vars
var (
	installPack string
	removePack  string
	searchTerm  string
	updatePack  string
)

func init() {
	flaggy.SetName("dfpm")
	flaggy.SetDescription("Device Family Pack Manager")
	flaggy.DefaultParser.AdditionalHelpPrepend = "Download, install, and search Microchip DFPs."
	flaggy.DefaultParser.AdditionalHelpAppend = "Call each subcommand with the help flag for more inforation about the usage."
	flaggy.SetVersion(Dfpm_Version)
	flaggy.DisableCompletion()

	flaggy.String(&packDir, "p", "pack-dir", "Specify pack directory to use.")

	cmdInstall = flaggy.NewSubcommand("install")
	cmdInstall.Description = "Install a given device family pack."
	cmdInstall.AddPositionalValue(&installPack, "name", 1, true, "Pack to install.")

	cmdRemove = flaggy.NewSubcommand("remove")
	cmdRemove.Description = "Remove a given device family pack."
	cmdRemove.AddPositionalValue(&removePack, "name", 1, true, "Pack to remove.")
	cmdRemove.Bool(&yesRemove, "y", "yes", "Remove pack(s) without confirmation.")

	cmdList = flaggy.NewSubcommand("list")
	cmdList.Description = "List device family packs."
	cmdList.Bool(&listAll, "", "all", "List all available packs, even ones not installed.")

	cmdSearch = flaggy.NewSubcommand("search")
	cmdSearch.Description = "Search for a given term in device packs."
	cmdSearch.AddPositionalValue(&searchTerm, "term", 1, true, "Term to search for.")
	cmdSearch.Bool(&downloadPacks, "d", "download", "Automatically download found packs.")
	cmdSearch.Bool(&yesDownload, "y", "yes", "Download packs without confirmation.")
	cmdSearch.Bool(&localSearch, "l", "local", "Only search locally installed packs.")

	cmdUpdate = flaggy.NewSubcommand("update")
	cmdUpdate.Description = "Update installed pack(s)."
	cmdUpdate.AddPositionalValue(&updatePack, "name", 1, false, "Pack to update. If not given, update all packs.")
	cmdUpdate.Bool(&showUpdates, "s", "show", "Show pending updates but do not update packs.")
	cmdUpdate.Bool(&keepPacks, "k", "keep", "Keep outdated packs.")

	cmdXc8 = flaggy.NewSubcommand("xc8")
	cmdXc8.Description = "Launch xc8-cc with -mdfp argument based on -mcpu."
	cmdXc8.AdditionalHelpPrepend = "If device is in multiple packs the command will fail."
	cmdXc8.AdditionalHelpPrepend += "\nPack will be automatically installed if not available."
	cmdXc8.AdditionalHelpAppend = "Additonal help append"

	flaggy.AttachSubcommand(cmdInstall, 1)
	flaggy.AttachSubcommand(cmdRemove, 1)
	flaggy.AttachSubcommand(cmdList, 1)
	flaggy.AttachSubcommand(cmdSearch, 1)
	flaggy.AttachSubcommand(cmdUpdate, 1)
	flaggy.AttachSubcommand(cmdXc8, 1)

	if len(os.Args) > 1 {
		// Loop through the first 3 arguments and check if the xc8 command is used
		for _, arg := range os.Args[1:min(4, len(os.Args)-1)] {
			if arg == "xc8" {
				// Disable errors on unexpected arguments
				flaggy.DefaultParser.ShowHelpOnUnexpected = false
				break
			}
		}
	}

	packDir = ExpandPath(packDir)

	flaggy.Parse()
}

func main() {
	var err error

	if cmdInstall.Used {
		err = Dfpm_install(packDir, installPack)
	} else if cmdRemove.Used {
		err = Dfpm_remove(packDir, removePack, yesRemove)
	} else if cmdList.Used {
		err = Dfpm_list(packDir, listAll)
	} else if cmdSearch.Used {
		err = Dfpm_search(packDir, searchTerm, downloadPacks, yesDownload, localSearch)
	} else if cmdUpdate.Used {
		err = Dfpm_update(packDir, updatePack, showUpdates, keepPacks)
	} else if cmdXc8.Used {
		var args []string
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "--pack-dir" ||
				os.Args[i] == "-pack-dir" ||
				os.Args[i] == "--p" ||
				os.Args[i] == "-p" {
				i++ // We dont save the flag, and skip over the next argument too
				continue
			} else if strings.HasPrefix(os.Args[i], "--pack-dir=") ||
				strings.HasPrefix(os.Args[i], "-pack-dir=") ||
				strings.HasPrefix(os.Args[i], "--p=") ||
				strings.HasPrefix(os.Args[i], "-p=") {
				continue // Don't save this flag
			} else if os.Args[i] == "xc8" {
				continue // Don't save the subcommand
			}
			args = append(args, os.Args[i])
		}
		err = Dfpm_xc8(packDir, args)
	} else {
		flaggy.ShowHelp("")
	}

	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}
