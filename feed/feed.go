package feed

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"os/exec"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mmcdole/gofeed"
)

const confdir string = "/.config/sreader"
const datadir string = "/.local/share/sreader"
const idxfile string = datadir + "/index"
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

func Init() []*gofeed.Feed {
	/* set configuration stuff */
	urlsfile := os.Getenv("HOME") + confdir + "/urls"

	/* this won't do anything if the files exist already */
	os.MkdirAll(os.Getenv("HOME") + confdir, os.ModePerm)
	os.MkdirAll(os.Getenv("HOME") + datadir, os.ModePerm)

	_, err := os.Stat(urlsfile)
    if os.IsNotExist(err) {
		file, err := os.Create(urlsfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	dat, err := ioutil.ReadFile(urlsfile)

	if err != nil {
		panic(err)
	}

	urls = strings.Split(string(dat), "\n")

func LoadFeeds() []*gofeed.Feed {
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

		if resp.StatusCode == http.StatusOK {
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				panic(err)
			}
		} else {
			panic("Failed to download feed \"" + url + "\": " + resp.Status)
		}
	}
}
