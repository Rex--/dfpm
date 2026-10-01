package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// Source - https://stackoverflow.com/a/24792688
// Posted by Astockwell, modified by community. See post 'Timeline' for change history
// Retrieved 2026-07-17, License - CC BY-SA 4.0
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExecuteCommand(cmd string, args []string) error {
	exe := exec.Command(cmd, args...)
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr

	return exe.Run()
}

func ExpandPath(path string) string {
	usr, _ := user.Current()
	hdir := usr.HomeDir

	if path == "~" {
		return hdir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(hdir, path[2:])
	}

	return path
}

var (
	ErrInvalidPackName = errors.New("invalid pack name")
	ErrPackNotFound    = errors.New("pack does not exist")
)

// This function takes various dfp names and returns a fully formed pack name.
// Acceptable names include: (given pack name Microchip.PIC16F1xxxx_DFP.1.32.471)
// - microchip.pic16f1xxxx_dfp.1.32.471.atpack
// - microchip.pic16f1xxxx_dfp.1.32.471
// - microchip.pic16f1xxxx_dfp
// - pic16f1xxxx_dfp
// - pic16f1xxxx
// - pic16f1
// NOTE: Names are not case sensitive. If no version is given, the latest version will be assumed.
//
// Always returns a name in format Microchip.<name>.<MajorVer>.<MinorVer>.<DotVer>
// e.g. 'Microchip.PIC16F1xxxx_DFP.1.32.471'
// Returns ErrInvalidPackName when name cannot be determined.
// Returns ErrPackNotFound if given pack doesn't exist.
func NormalizePackName(name string) (string, error) {
	// microchip.pic16f1xxxx_dfp.1.32.471.atpack -> pic16f1xxxx_dfp.1.32.471
	name = strings.TrimPrefix(strings.TrimSuffix(name, ".atpack"), "Microchip.")

	basename, version, hasVer := strings.Cut(name, ".")

	webPacks, err := FindWeb()
	if err != nil {
		return "", err
	}

	var packName string
	for _, wp := range webPacks {
		if strings.HasPrefix(strings.ToLower(wp.Name), strings.ToLower(basename)) {
			// We've matched a packname. If already matched one, return ErrInvalidPackName
			if packName != "" {
				return "", ErrInvalidPackName
			}

			if hasVer {
				// Verify version (TODO: look up available versions)
				// Check if there is exactly 2 periods that are separated by at least one number
				verList := strings.Split(version, ".")
				if len(verList) != 3 {
					return "", ErrInvalidPackName
				}
				for _, n := range verList {
					if len(n) == 0 {
						return "", ErrInvalidPackName
					}
					for _, c := range n {
						if c < '0' || c > '9' {
							return "", ErrInvalidPackName
						}
					}
				}
				packName = strings.Join([]string{"Microchip", wp.Name, version}, ".")
			} else {
				packName = strings.Join([]string{"Microchip", wp.Name, wp.Version}, ".")
			}
		}
	}

	if packName == "" {
		return "", ErrPackNotFound
	}

	return packName, nil
}
