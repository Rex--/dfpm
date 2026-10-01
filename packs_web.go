package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type WebPack struct {
	Name    string
	Version string
	Link    string
	Search  string
}

var cachedWebPacks []*WebPack

var DownloadURL = "https://packs.download.microchip.com/"

func InstallPack(pack_name, install_dir string) (err error) {
	pack_file, err := downloadPack(pack_name)
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

func GetPackIndex() (packs []*WebPack, err error) {

	// Return cached web packs if weve already done this
	if cachedWebPacks != nil {
		return cachedWebPacks, nil
	}

	resp, err := http.Get(DownloadURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	d := xml.NewDecoder(resp.Body)

	start, err := skipToData(d)
	if err != nil {
		return nil, err
	}

	var pl packList
	err = d.DecodeElement(&pl, &start)
	if err != nil {
		return nil, err
	}

	for _, wp := range pl.Packs {
		packName := wp.Details.Href
		nameParts := strings.SplitN(wp.Details.Href, ".", 3)
		if len(nameParts) == 3 {
			packName = nameParts[1]
		}
		packVer := strings.TrimPrefix(wp.Details.Version, "v")

		pack := WebPack{
			Name:    packName,
			Version: packVer,
			Link:    wp.Details.Href,
			Search:  wp.Search,
		}
		packs = append(packs, &pack)
	}

	// save the generated list so we dont have to hit the url and parse again
	cachedWebPacks = packs

	return
}

func FindWeb() (packs []Pack, err error) {

	wpacks, err := GetPackIndex()
	if err != nil {
		return nil, err
	}

	for _, webPack := range wpacks {
		pack := Pack{
			Name:    webPack.Name,
			Version: webPack.Version,
		}
		packs = append(packs, pack)
	}

	return
}

func SearchWeb(term string) (packs []Pack, err error) {
	webPacks, err := GetPackIndex()
	if err != nil {
		return nil, err
	}

	for _, wp := range webPacks {
		if strings.Contains(strings.ToLower(wp.Search), strings.ToLower(term)) {
			pack := Pack{
				Name:    wp.Name,
				Version: wp.Version,
			}
			packs = append(packs, pack)
		}
	}

	return
}

type packList struct {
	XMLName xml.Name  `xml:"div"`
	Packs   []article `xml:"article"`
}

type article struct {
	Search  string  `xml:"data-search,attr"`
	Details details `xml:"details>summary>a"`
}

type details struct {
	Href    string `xml:"href,attr"`
	Version string `xml:",innerxml"`
}

func skipToData(d *xml.Decoder) (xml.StartElement, error) {
	for {
		t, err := d.Token()
		if err != nil {
			return xml.StartElement{}, err
		}

		switch v := t.(type) {
		case xml.StartElement:
			if v.Name.Local == "div" {
				for _, attr := range v.Attr {
					if attr.Name.Local == "class" && attr.Value == "pack-list" {
						return v.Copy(), nil
					}
				}
			}
		}
	}
}

func downloadPack(pack_name string) (pack_file string, err error) {

	url, err := url.JoinPath(DownloadURL, pack_name)
	if err != nil {
		return
	}

	file, err := os.CreateTemp("", pack_name)
	if err != nil {
		return
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("downloading", url)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}

	return file.Name(), nil

}
