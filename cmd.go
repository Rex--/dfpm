package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"
)

func Dfpm_install(packDir, name string) error {

	normName, err := NormalizePackName(name)
	if err != nil {
		return err
	}

	// Check to see if pack is already installed
	nameParts := strings.SplitN(normName, ".", 3)
	packs, err := FindLocal(packDir)
	if err != nil {
		return err
	}
	for _, p := range packs {
		if p.Name == nameParts[1] {
			if p.Version == nameParts[2] {
				fmt.Println("pack already installed, nothing to do")
				return nil
			} else {
				fmt.Printf("pack found with version %s, installing new pack with version %s\n", p.Version, nameParts[2])
				// fmt.Printf("remove old pack with `%s remove %s.%s`\n", os.Args[0], p.Name, p.Version)
			}
		}
	}

	err = InstallPack(normName+".atpack", packDir)

	return err
}

func Dfpm_remove(packDir, name string, yes bool) error {
	packs, err := SearchLocal(packDir, name)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d packs matching '%s'\n", len(packs), name)
	if len(packs) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)
		for _, p := range packs {
			fmt.Fprintf(w, "  %s\tv%s\n", p.Name, p.Version)
		}
		w.Flush()

		var resp string
		if !yes {
			fmt.Printf("Remove %d packs? [Y/n]: ", len(packs))
			fmt.Scanln(&resp)
		}

		if resp == "" || strings.HasPrefix(strings.ToLower(resp), "y") {
			for _, p := range packs {
				fmt.Printf("Removing %s\n", p.Path)
				err := os.RemoveAll(p.Path)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func Dfpm_list(packDir string, all bool) error {
	packs, err := FindLocal(packDir)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)
	// fmt.Fprintf(w, "Installed Packs:\n")
	if len(packs) == 0 {
		fmt.Fprintf(w, "No packs found in %s\n", packDir)
	} else {
		for _, pack := range packs {
			fmt.Fprintf(w, "%s\tv%s\n", pack.Name, pack.Version)
		}
	}
	if all {
		packs, err = FindWeb()
		if err != nil {
			return err
		}
		fmt.Fprintln(w, "\nAvailable Packs:")
		for _, pack := range packs {
			fmt.Fprintf(w, "%s\tv%s\n", pack.Name, pack.Version)
		}
	}
	w.Flush()

	return nil
}

func Dfpm_search(packDir string, term string, download bool, yes bool, local bool) error {

	localPacks, err := SearchLocal(packDir, term)
	if err != nil {
		return err
	}
	var localNames []string
	for _, pack := range localPacks {
		localNames = append(localNames, pack.Name)
	}

	var packs []Pack
	if local {
		packs, err = SearchLocal(packDir, term)
	} else {
		packs, err = SearchWeb(term)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Found '%s' in %d packs\n", term, len(packs))
	anyInstalled := false
	var toInstall []string
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)
	for _, p := range packs {
		if slices.Contains(localNames, p.Name) {
			anyInstalled = true
			if local { // Show version with local packages (useful when we have multiple versions)
				fmt.Fprintf(w, "* %s\tv%s\n", p.Name, p.Version)
			} else { // No version because each pack is only listed once
				fmt.Fprintf(w, "* %s\n", p.Name)
			}
		} else {
			fmt.Fprintf(w, "  %s\n", p.Name)
			if download {
				packName := strings.Join([]string{"Microchip", p.Name, p.Version, "atpack"}, ".")
				toInstall = append(toInstall, packName)
				// InstallPack(packName, packDir)
			}
		}
	}
	w.Flush()
	if len(packs) > 0 && anyInstalled {
		fmt.Println(" * Installed")
	}
	if download && len(toInstall) > 0 {
		var resp string
		if !yes {
			fmt.Printf("Install %d packs? [Y/n]: ", len(toInstall))
			fmt.Scanln(&resp)
		}

		if yes || resp == "" || strings.HasPrefix(strings.ToLower(resp), "y") {
			for _, p := range toInstall {
				err = InstallPack(p, packDir)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func Dfpm_update(packDir string, name string, show bool, keepPacks bool) error {
	localPacks, err := FindLocal(packDir)
	if err != nil {
		return err
	}

	webPacks, err := FindWeb()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)

	if name == "" { // Update all packs if not specified
		upToDate := true
		for _, pack := range localPacks {
			for _, wpack := range webPacks {
				if pack.Name == wpack.Name && pack.Version != wpack.Version {
					// We decide if its out of date by checking if its different. Hopefully microchip doesn't release an older pack than what we have installed.
					upToDate = false
					fmt.Fprintf(w, "%s\tv%s ->\t%s\n", pack.Name, pack.Version, wpack.Version)
					if !show {
						w.Flush()
						packName := strings.Join([]string{"Microchip", pack.Name, wpack.Version, "atpack"}, ".")
						InstallPack(packName, packDir)
						if !keepPacks {
							fmt.Printf("removing old pack %s v%s\n", pack.Name, pack.Version)
							err = os.RemoveAll(pack.Path)
							if err != nil {
								return err
							}
						}
					}
				} else if show && pack.Name == wpack.Name {
					fmt.Fprintf(w, "%s\tv%s\tlatest\n", pack.Name, pack.Version)
				}
			}
		}
		if !show && upToDate {
			fmt.Fprintln(w, "All packs up to date.")
		}
	} else { // Update specific pack
		normName, err := NormalizePackName(name)
		if err != nil {
			return err
		}
		nameParts := strings.SplitN(normName, ".", 3)
		for _, pack := range localPacks {
			if pack.Name == nameParts[1] {
				for _, wpack := range webPacks {
					if pack.Name == wpack.Name && pack.Version != wpack.Version {
						// We decide if its out of date by checking if its different. Hopefully microchip doesn't release an older pack than what we have installed.
						fmt.Fprintf(w, "%s\tv%s ->\t%s\n", pack.Name, pack.Version, wpack.Version)
						if !show {
							w.Flush()
							// NOTE: use version from normalized name here. If not specified it will be the latest, else we can up(or down)grade to a specific version
							packName := strings.Join([]string{"Microchip", pack.Name, nameParts[2], "atpack"}, ".")
							InstallPack(packName, packDir)
							if !keepPacks {
								fmt.Printf("removing old pack %s v%s\n", pack.Name, pack.Version)
								err = os.RemoveAll(pack.Path)
								if err != nil {
									return err
								}
							}
						}
					} else if pack.Name == wpack.Name {
						fmt.Fprintf(w, "%s\tv%s\tlatest\n", pack.Name, pack.Version)
					}
				}
			}
		}
	}
	w.Flush()

	return nil
}

func Dfpm_xc8(packDir string, args []string) error {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-mcpu=") {
			mcpu := strings.TrimPrefix(arg, "-mcpu=")
			packs, err := SearchLocal(packDir, mcpu)
			if err != nil {
				return err
			} else if len(packs) > 1 {
				return fmt.Errorf("device %s is in multiple packs", mcpu)
			}

			// Couldn't find local pack with device
			if len(packs) == 0 {
				packs, err = SearchWeb(mcpu)
				if err != nil {
					return err
				} else if len(packs) == 0 {
					return fmt.Errorf("device %s not found", mcpu)
				} else if len(packs) > 1 {
					return fmt.Errorf("device %s is in multiple packs", mcpu)
				}

				fmt.Printf("Found device %s in pack %s, installing...\n", mcpu, packs[0].Name)

				// Install pack
				packName := strings.Join([]string{"Microchip", packs[0].Name, packs[0].Version, "atpack"}, ".")
				err = InstallPack(packName, packDir)
				if err != nil {
					return err
				}
				packs[0].Local = true
				packs[0].Path = filepath.Join(packDir, strings.TrimSuffix(packName, ".atpack"))
			}

			mdfp := filepath.Join(packs[0].Path, "xc8")
			args = slices.Insert(args, 0, "-mdfp="+mdfp)
			break
		}
	}

	err := ExecuteCommand("xc8-cc", args)
	if err != nil {
		return err
	}

	return nil
}
