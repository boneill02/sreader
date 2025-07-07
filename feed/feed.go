package feed

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"html"
	"io"
	"log"
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

// Initialize the backend.
// This function does the following:
// - calls InitDB to initialize the database,
// - loads user-specified URLs
func Init() {
	InitDB()

	/* set configuration stuff */
	urlsfile := config.Config.ConfDir + "/urls"
	_, err := os.Stat(urlsfile)
	if os.IsNotExist(err) {
		file, err := os.Create(urlsfile)
		if err != nil {
			log.Fatalln("Failed to create URLs file", err.Error())
		}
		defer file.Close()
	}
	dat, err := os.ReadFile(urlsfile)
	if err != nil {
		log.Fatalln("Failed to read URLs file", err.Error())
	}
	urls = strings.Split(string(dat), "\n")
}

// Open URL in web browser
func OpenInBrowser(url string, browser string) {
	cmd := exec.Command(browser, url)
	cmd.Start()
}

// Open URL in media player
func OpenInPlayer(url string, player string) {
	cmd := exec.Command("setsid", "nohup", player, url)
	cmd.Start()
}

// Sync feeds
// This function asynchronously GETs feeds, using the last_updated field
// in the database to only grab/update feeds that were updated since the last sync.
// The new feed contents are then stored in the database.
func Sync() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	// Listen for OS signals to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start workers
	log.Println("Getting feeds...")
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

	log.Println("Updating DB...")
	feed_contents := loadRSSFeeds()
	for _, f := range feed_contents {
		if f != nil {
			if id, err := AddFeed(f); err != nil {
				log.Println("Error adding feed:", err.Error())
			} else {
				MarkUpdated(id)
			}
		}
	}
	log.Println("Done.")
}

// Unescape HTML entities and convert to ASCII
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

// Called by loadRSSFeeds. Parse feed from temporary file grabbed by syncWorkers and remove the file.
func loadRSSFeed(url string) *gofeed.Feed {
	urlsum := sha1.Sum([]byte(url))
	filename := config.Config.DataDir + "/" + hex.EncodeToString(urlsum[:]) + ".tmp"
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Failed to open temporary file:", err.Error())
		return nil
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)

	if err != nil {
		log.Println("Failed to parse feed (possibly wrong URL or badly formatted XML?)")
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

// Parse all downloaded RSS feeds
func loadRSSFeeds() []*gofeed.Feed {
	var feeds []*gofeed.Feed

	for _, url := range urls {
		if len(url) > 0 {
			feeds = append(feeds, loadRSSFeed(url))
		}
	}

	return feeds
}

// Single GET request worker
func syncWorker(url string, modTime string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	// Get file name for URL
	urlsum := sha1.Sum([]byte(url))
	filename := config.Config.DataDir + "/" + hex.EncodeToString(urlsum[:]) + ".tmp"

	// Create request to fetch the feed
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Failed to create request for URL:", url, "Error:", err)
		return
	}

	req.Header.Set("User-Agent", "sreader/1.0")
	if modTime != "" {
		req.Header.Set("If-Modified-Since", modTime)
	}

	// Do GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Failed to fetch feed:", url, "Error:", err)
	}

	// Create the temporary file
	out, err := os.Create(filename)
	if err != nil {
		log.Fatalln("Failed to create temporary file:", err)
	}

	// Copy response body to the temporary file
	if resp.StatusCode == http.StatusOK {
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}
	} else if resp.StatusCode != http.StatusNotModified {
		log.Println("Failed to download feed \"" + url + "\": " + resp.Status)
	}

	out.Close()

	select {
	case <-ctx.Done():
		// Context was cancelled, clean up and exit
		log.Println("Sync cancelled for URL:", url)
		os.Remove(filename) // Clean up temporary file
		return
	default:
		return
	}
}
