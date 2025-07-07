package feed

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"html"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/bmoneill/sreader/config"
	"github.com/mmcdole/gofeed"
)

var urls []string

func Init() {
	InitDB()
	/* set configuration stuff */
	urlsfile := config.Config.ConfDir + "/urls"
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
}

/**
 * Open feed in web browser
 * Uses the BROWSER environment variable to determine which browser to use.
 * If BROWSER is not set, it will not open the URL.
 */
func OpenInBrowser(url string, browser string) {
	cmd := exec.Command(browser, url)
	cmd.Start()
}

/**
 * Open feed in video player
 */
func OpenInPlayer(url string, player string) {
	cmd := exec.Command("setsid", "nohup", player, url)
	cmd.Start()
}

/**
 * Sync all feeds.
 */
func Sync() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	// Listen for OS signals to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start workers
	println("Getting feeds...")
	for i := range urls {
		url := urls[i]
		if len(url) < 1 {
			continue
		}

		feed := GetFeedByURL(url)
		wg.Add(1)
		if feed != nil {
			go syncWorker(url, feed.LastUpdated, &wg, ctx)
		} else {
			go syncWorker(url, "", &wg, ctx)
		}
	}

	go func() {
		<-sigChan
		cancel() // Cancel context when signal received
	}()

	wg.Wait()

	println("Updating DB...")
	feed_contents := loadRSSFeeds()
	for _, f := range feed_contents {
		if f != nil {
			if id, err := AddFeed(f); err != nil {
				println("Error adding feed:", err.Error())
			} else {
				MarkUpdated(id)
			}
		}
	}
}

/**
 * Unescape HTML entities and convert to ASCII
 */
func formatHTMLString(s string) string {
	s = html.UnescapeString(s)

	ascii := make([]rune, 0, len(s))
	for _, r := range s {
		if r < 128 {
			ascii = append(ascii, r)
		} else {
			ascii = append(ascii, ' ')
		}
	}
	return string(ascii)
}

/**
 * Parse feed from data directory
 */
func loadRSSFeed(url string) *gofeed.Feed {
	urlsum := sha1.Sum([]byte(url))
	filename := config.Config.DataDir + "/" + hex.EncodeToString(urlsum[:]) + ".tmp"
	file, err := os.Open(filename)
	println("Loading feed from file:", filename)

	if err != nil {
		// try to sync feed if it doesn't exist
		println("Feed file not found for URL:", url, "Error:", err)
		Sync()
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)

	if err != nil {
		println("Failed to parse feed from file:", filename, "Error:", err.Error())
		return nil
	}

	// Unescape HTML entities and convert to ASCII
	feed.Description = formatHTMLString(feed.Description)
	feed.Title = formatHTMLString(feed.Title)

	for _, item := range feed.Items {
		item.Title = formatHTMLString(item.Title)
		item.Description = formatHTMLString(item.Description)
		item.Content = formatHTMLString(item.Content)
	}

	os.Remove(filename) // Clean up temporary file
	return feed
}

func loadRSSFeeds() []*gofeed.Feed {
	var feeds []*gofeed.Feed

	for _, url := range urls {
		if len(url) > 0 {
			feeds = append(feeds, loadRSSFeed(url))
		}
	}

	return feeds
}

func syncWorker(url string, modTime string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	print(modTime)

	// Get file name for URL
	urlsum := sha1.Sum([]byte(url))
	filename := config.Config.DataDir + "/" + hex.EncodeToString(urlsum[:]) + ".tmp"

	// Create request to fetch the feed
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		println("Failed to create request for URL:", url, "Error:", err)
		return
	}

	req.Header.Set("User-Agent", "sreader/1.0")
	if modTime != "" {
		req.Header.Set("If-Modified-Since", modTime)
	}

	// Do GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	// Create the temporary file
	out, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Copy response body to the temporary file
	if resp.StatusCode == http.StatusOK {
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}
	} else if resp.StatusCode != http.StatusNotModified {
		panic("Failed to download feed \"" + url + "\": " + resp.Status)
	}

	out.Close()

	select {
	case <-ctx.Done():
		// Context was cancelled, clean up and exit
		println("Sync cancelled for URL:", url)
		os.Remove(filename) // Clean up temporary file
		return
	default:
		return
	}
}
