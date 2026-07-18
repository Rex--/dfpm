package main

import (
	"archive/zip"
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

func ListSubDirectories(path string) (dirs []string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}

	return
}

func ExecuteCommand(cmd []string) error {
	exe := exec.Command(cmd[0], cmd[1:]...)
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
