package feed

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/mmcdole/gofeed"
)

const confdir string = "/.config/sreader"
const datadir string = "/.local/share/sreader"

var urls []string

/**
 * parse feed from data directory
 */
func GetFeed(url string) *gofeed.Feed {
	urlsum := sha1.Sum([]byte(url))
	file, err := os.Open(os.Getenv("HOME") + datadir + "/" + hex.EncodeToString(urlsum[:]))

	if err != nil {
		panic(err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)

	if err != nil {
		panic(err)
	}

	return feed
}

/* Get URLs from the urls file */
func GetUrls() []string {
	urlsfile := os.Getenv("HOME") + confdir + "/urls"
	_, err := os.Stat(urlsfile)
	if os.IsNotExist(err) {
		file, err := os.Create(urlsfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}
	dat, err := os.ReadFile(urlsfile)
	if err != nil {
		panic(err)
	}
	urls = strings.Split(string(dat), "\n")
	return urls
}

/* Create all necessary files and directories if they don't exist yet */
func CreateFiles() {
	urlsfile := os.Getenv("HOME") + confdir + "/urls"
	datadir := os.Getenv("HOME") + datadir

	// create urls file if it doesn't exist
	_, err := os.Stat(urlsfile)
	if os.IsNotExist(err) {
		file, err := os.Create(urlsfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	// create data directory if it doesn't exist
	os.MkdirAll(datadir, os.ModePerm)
}

func Init() []*gofeed.Feed {
	CreateFiles()
	urls := GetUrls()

	var feeds []*gofeed.Feed

	for _, url := range urls {
		if len(url) > 0 {
			feeds = append(feeds, GetFeed(url))
		}
	}

	return feeds
}

/* open feed in default browser */
func OpenInBrowser(url string) {
	browser := os.Getenv("BROWSER")
	if browser != "" {
		cmd := exec.Command(browser, url)
		cmd.Start()
	}
}

/* open feed in video player */
func OpenInPlayer(url string) {
	player := os.Getenv("PLAYER")

	if player == "" {
		player = "mpv" // default player
	}

	cmd := exec.Command("setsid", "nohup", player, url)
	cmd.Start()
}

/**
 * sync all feeds (download files)
 */
func Sync() {
	CreateFiles()
	GetUrls()

	for _, url := range urls {
		if len(url) < 1 {
			continue
		}
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		urlsum := sha1.Sum([]byte(url))
		filename := os.Getenv("HOME") + datadir + "/" + hex.EncodeToString(urlsum[:])
		out, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}
	}
}
