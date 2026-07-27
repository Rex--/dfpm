package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func Dfpm_install(args DFPMArgs) {
	var pack string
	if len(args.Args) == 1 {
		pack = args.Args[0]
		// println("installing pack:", pack)
	} else {
		fmt.Println("missing required argument: <pack_name>")
		os.Exit(1)
	}

	err := InstallPack(pack, args.PackDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func Dfpm_list(args DFPMArgs) {
	packs := ListSubDirectories(args.PackDir)
	for _, p := range packs {
		fmt.Println(filepath.Join(args.PackDir, p))
	}
}

func Dfpm_search(args DFPMArgs) {

	var search string
	if len(args.Args) == 1 {
		search = args.Args[0]
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

	packs_dl := ListSubDirectories(args.PackDir)

	var packs_2dl []string

	has := false
	fmt.Printf("Found '%s' in %d packs\n", search, len(packs))
	for _, p := range packs {
		dld := "  "
		if slices.Contains(packs_dl, strings.TrimSuffix(p, ".atpack")) {
			dld = "* "
			has = true
		} else {
			packs_2dl = append(packs_2dl, p)
		}
		fmt.Printf("  %s%s\n", dld, p)
	}
	if has {
		fmt.Println(" (* Downloaded)")
	}

	if args.Download {
		if len(packs_2dl) > 1 {
			fmt.Printf("\nDownload %d packs? [Y/n]: ", len(packs_2dl))
			var resp string
			fmt.Scanln(&resp)

			if resp == "" || strings.HasPrefix(strings.ToLower(resp), "y") {
				for _, p := range packs_2dl {
					err := InstallPack(p, args.PackDir)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
				}
			}
		}
	}

}

func dfpm_update() {

}

func Dfpm_xc8(args DFPMArgs) {
	dfpm := ""
	for _, a := range args.Args {
		if strings.HasPrefix(a, "-mcpu=") {
			mcpu := strings.TrimPrefix(a, "-mcpu=")
			f, err := GetPackIndex()
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}
			pack := ParseHtmlXml(f, mcpu)
			packs_dl := ListSubDirectories(args.PackDir)
			if len(pack) == 1 {
				pack_name := strings.TrimSuffix(pack[0], ".atpack")
				if slices.Contains(packs_dl, pack_name) {
					dfpm, _ = filepath.Abs(filepath.Join(args.PackDir, pack_name, "xc8"))
					args.Args = slices.Insert(args.Args, 0, "-mdfp="+dfpm)
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

	args.Args = slices.Insert(args.Args, 0, "xc8-cc")

	// fmt.Println(strings.Join(args.Args, " "))

	err := ExecuteCommand(args.Args)
	if err != nil {
		fmt.Println("dfpm-xc8:", err.Error())
		os.Exit(1)
	}
}
