package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var DownloadURL = "https://packs.download.microchip.com/"

func InstallPack(pack_name, install_dir string) (err error) {
	pack_file, err := DownloadPack(pack_name)
	if err != nil {
		return
	}

	pack_install, _ := filepath.Abs(filepath.Join(install_dir, strings.TrimSuffix(pack_name, ".atpack")))

	fmt.Println("unpacking", pack_name)
	err = Unzip(pack_file, pack_install)
	if err != nil {
		return
	}

	// Delete downloaded packfile
	err = os.Remove(pack_file)
	if err != nil {
		return
	}

	fmt.Println("installed", pack_install)

	return
}

func DownloadPack(pack_name string) (pack_file string, err error) {

	// filename := path.Join(InstallDir, "downloads", pack_name)
	url, err := url.JoinPath(DownloadURL, pack_name)
	if err != nil {
		return
	}

	// err = os.MkdirAll(path.Dir(filename), os.ModePerm)
	// if err != nil {
	// 	return
	// }

	file, err := os.CreateTemp("", pack_name)
	if err != nil {
		return
	}
	// file, err := os.Create(filename)
	// if err != nil {
	// 	return
	// }
	defer file.Close()
	// println(url, "->", file.Name())
	fmt.Println("downloading", url)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}

	return file.Name(), nil

}

func GetPackIndex() (r io.Reader, err error) {

	resp, err := http.Get(DownloadURL)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return resp.Body, nil
}
