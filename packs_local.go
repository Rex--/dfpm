package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Pack struct {
	Name    string // Name of pack e.g. PIC16F1xxxx_DFP
	Version string // Pack version e.g. 1.32.471

	Local bool   // Is pack local? (has been downloaded and installed)
	Path  string // Path to downloaded pack if Local==true
}

func FindLocal(packDir string) (packs []Pack, err error) {

	packPaths, err := getPacksInDir(packDir)
	if err != nil {
		return nil, err
	}

	for _, path := range packPaths {
		pathName := filepath.Base(path)
		packList := strings.SplitN(pathName, ".", 3)
		packName := pathName
		packVersion := "unknown"
		if len(packList) == 3 {
			packName = packList[1]
			packVersion = packList[2]
		}
		pack := Pack{
			Name:    packName,
			Version: packVersion,
			Local:   true,
			Path:    path,
		}
		packs = append(packs, pack)
	}

	return
}

func SearchLocal(packDir string, term string) (packs []Pack, err error) {
	term = strings.ToLower(term)

	packPaths, err := getPacksInDir(packDir)
	if err != nil {
		return nil, err
	}

packPathsLoop:
	for _, path := range packPaths {
		// Look for pack names that contain the search term
		if strings.Contains(strings.ToLower(filepath.Base(path)), term) {
			pathName := filepath.Base(path)
			packList := strings.SplitN(pathName, ".", 3)
			packName := pathName
			packVersion := "unknown"
			if len(packList) == 3 {
				packName = packList[1]
				packVersion = packList[2]
			}
			pack := Pack{
				Name:    packName,
				Version: packVersion,
				Local:   true,
				Path:    path,
			}
			packs = append(packs, pack)
		} else {
			// Look for device files that contain the search term
			files, err := os.ReadDir(filepath.Join(path, "edc"))
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				if strings.Contains(strings.ToLower(file.Name()), term) {
					pathName := filepath.Base(path)
					packList := strings.SplitN(pathName, ".", 3)
					packName := pathName
					packVersion := "unknown"
					if len(packList) == 3 {
						packName = packList[1]
						packVersion = packList[2]
					}
					pack := Pack{
						Name:    packName,
						Version: packVersion,
						Local:   true,
						Path:    path,
					}
					packs = append(packs, pack)
					continue packPathsLoop // Skip the rest of this pack after adding
				}
			}
		}
	}

	return
}

// This function returns a list of all DFPs in a directory.
// Searches recursively and returns the full path for each.
func getPacksInDir(path string) (packPaths []string, err error) {

	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipAll
		}
		if d.IsDir() {
			// Check for a Microchip.xxxx.pdsc file in directory
			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".pdsc" && strings.HasPrefix(file.Name(), "Microchip.") {
					packPaths = append(packPaths, path)
					return filepath.SkipDir
				}
			}
		}

		return nil
	})

	return
}
