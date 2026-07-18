package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var InstallDir = "~/.mchp_packs"

func main() {
	InstallDir = ExpandPath(InstallDir)

	if len(os.Args) == 1 {
		usage()
	} else {
		switch os.Args[1] {
		case "install":
			dfpm_install()
		case "list":
			dfpm_list()
		case "search":
			dfpm_search()
		case "update":
			dfpm_update()
		case "xc8":
			dfpm_xc8()
		default:
			usage()
		}
	}

}

func usage() {
	usage :=
		`Usage: dfpm [install | list | search ] [arg]
       dfpm install <pack>      install the given device pack
       dfpm list                list installed device packs
       dfpm search <string>     search for <string> in pack descriptions

       dfpm xc8 [-mcpu=<cpuid>] [args ...]
                                search for the given <cpuid> in pack descriptions
                                if found, launch 'xc8-cc' with the --mdfp option`
	fmt.Println(usage)
}

func dfpm_install() {
	var pack string
	if len(os.Args) > 2 {
		pack = os.Args[2]
		println("installing pack:", pack)
	} else {
		fmt.Println("missing required argument: <pack_name>")
		os.Exit(1)
	}

	err := InstallPack(pack)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func dfpm_list() {
	// installpath, _ := filepath.Abs(InstallDir)
	packs := ListSubDirectories(InstallDir)
	for _, p := range packs {
		fmt.Println(filepath.Join(InstallDir, p))
	}
}

func dfpm_search() {

	var search string
	if len(os.Args) > 2 {
		search = os.Args[2]
	} else {
		fmt.Println("missing required argument: <string>")
		os.Exit(1)
	}

	f, err := GetPackIndex()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	packs := ParseHtmlXml(f, search)

	packs_dl := ListSubDirectories(InstallDir)

	has := false
	fmt.Printf("Found '%s' in %d packs\n", search, len(packs))
	for _, p := range packs {
		dld := "  "
		if slices.Contains(packs_dl, strings.TrimSuffix(p, ".atpack")) {
			dld = "* "
			has = true
		}
		fmt.Printf("  %s%s\n", dld, p)
	}
	if has {
		fmt.Println(" (* Downloaded)")
	}
}

func dfpm_update() {

}

func dfpm_xc8() {
	args := os.Args[2:]
	dfpm := ""
	for _, a := range args {
		if strings.HasPrefix(a, "-mcpu=") {
			mcpu := strings.TrimPrefix(a, "-mcpu=")
			f, err := GetPackIndex()
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}
			pack := ParseHtmlXml(f, mcpu)
			packs_dl := ListSubDirectories(InstallDir)
			if len(pack) == 1 {
				pack_dir := strings.TrimSuffix(pack[0], ".atpack")
				if slices.Contains(packs_dl, pack_dir) {
					dfpm, _ = filepath.Abs(filepath.Join(InstallDir, pack_dir, "xc8"))
					args = slices.Insert(args, 0, "-mdfp="+dfpm)
					break
				} else {
					println("need to download pack:", pack[0])
					return
				}
			} else if len(pack) > 1 {
				println("mcpu matches multiple packs")
				return
			} else {
				println("unknown mcpu:", mcpu)
				return
			}
		}
	}

	args = slices.Insert(args, 0, "xc8-cc")

	// fmt.Println(strings.Join(args, " "))

	err := ExecuteCommand(args)
	if err != nil {
		fmt.Println("dfpm-xc8:", err.Error())
		os.Exit(1)
	}
}
