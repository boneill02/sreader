package feed

import (
	"database/sql"
	"log"

	"github.com/bmoneill/sreader/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
)

type Entry struct {
	ID            int64  `json:"id"`
	FeedID        int64  `json:"feed_id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	DatePublished string `json:"date_published"`
	Read          bool   `json:"read"`
}

type Feed struct {
	ID          int64    `json:"id"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	LastUpdated string   `json:"last_updated"`
	Entries     []*Entry `json:"entries,omitempty"`
}

var conn *sql.DB

// Initialize SQLite database connection and creates the necessary tables if they do not exist.
func InitDB() {
	log.Println("Initializing database...")
	// Initialize the SQLite database connection
	var err error
	conn, err = sql.Open("sqlite3", config.Config.DBFile)
	if err != nil {
		log.Fatalln("Failed to initialize database.", err.Error())
	}

	// Create the tables if they do not exist
	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT,
		title TEXT,
		description TEXT,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		log.Fatalln("Error creating feeds table:", err.Error())
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		feed_id INTEGER,
		url TEXT,
		title TEXT NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '',
		content TEXT NOT NULL DEFAULT '',
		date_published DATETIME,
		read INTEGER DEFAULT 0,
		FOREIGN KEY(feed_id) REFERENCES feeds(id)
	)`)

	if err != nil {
		log.Fatalln("Error creating entries table:", err.Error())
	}
}

// Adds a feed to the database.
// If the feed already exists, it adds any new entries.
// If the feed does not exist, it inserts a new feed and its entries.
func AddFeed(feed *gofeed.Feed) (int64, error) {
	var exists bool
	var id int64
	var stmt *sql.Stmt
	var res sql.Result
	var err error

	// Check if the feed already exists
	err = conn.QueryRow("SELECT EXISTS(SELECT 1 FROM feeds WHERE url = ?)", feed.Link).Scan(&exists)
	if err != nil {
		log.Println("Error checking if feed exists:", err.Error())
		return 0, err
	}

	if exists {
		// Get ID of existing feed
		log.Println("Feed already exists in DB, updating:", feed.Link)
		stmt, err = conn.Prepare("SELECT id FROM feeds WHERE url = ?")
		if err != nil {
			log.Println("Error preparing statement:", err.Error())
			return 0, err
		}
		defer stmt.Close()
		err = stmt.QueryRow(feed.Link).Scan(&id)
		if err != nil {
			log.Println("Error querying feed ID:", err.Error())
			return 0, err
		}
	} else {
		// Insert new feed into the database
		log.Println("Adding new feed to DB: ", feed.Link)
		stmt, err = conn.Prepare("INSERT INTO feeds (url, title, description) VALUES (?, ?, ?)")
		if err != nil {
			log.Println("Error preparing statement:", err.Error())
			return 0, err
		}
		defer stmt.Close()
		res, err = stmt.Exec(feed.Link, feed.Title, feed.Description)

		if err != nil {
			log.Println("Error inserting feed:", err.Error())
			return 0, err
		}
		id, _ = res.LastInsertId()
	}

	// Add entries
	for _, item := range feed.Items {
		err = AddEntry(id, item.Link, item.Title, item.Description, item.PublishedParsed.UTC().Format("Tue, 15 Nov 1994 12:45:26 GMT"))
		if err != nil {
			log.Println("Error adding entry:", err.Error())
			return 0, err
		}
	}

	log.Println(feed.Title, "added/updated successfully,", len(feed.Items), "entries.")
	return id, err
}

// Adds an entry to the database if it does not already exist.
func AddEntry(feedID int64, url, title, description string, datePublished string) error {
	// Check if the entry already exists (by feed_id and date_published)
	var exists bool
	err := conn.QueryRow("SELECT EXISTS(SELECT 1 FROM entries WHERE feed_id = ? AND date_published = ?)", feedID, datePublished).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // Entry already exists, do nothing
	}

	// Insert new entry into the database
	stmt, err := conn.Prepare("INSERT INTO entries (feed_id, url, title, description, date_published) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(feedID, url, title, description, datePublished)
	return err
}

// Get all entries for feed with feedID
func GetEntries(feedID int) []*Entry {
	// Retrieve entries for a specific feed
	rows, err := conn.Query("SELECT id, url, title, description, date_published, read, content FROM entries WHERE feed_id = ?", feedID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var entries []*Entry

	// Scan rows into Entry structs
	for rows.Next() {
		var (
			id            int
			url           string
			title         string
			description   string
			content       string
			datePublished string
			read          int
		)

		err := rows.Scan(&id, &url, &title, &description, &datePublished, &read, &content)
		if err != nil {
			log.Println("Error scanning entry:", err.Error())
			return nil
		}

		entry := &Entry{
			ID:            int64(id),
			FeedID:        int64(feedID),
			URL:           url,
			Title:         title,
			Description:   description,
			Content:       content,
			DatePublished: datePublished,
			Read:          read == 1,
		}
		entries = append(entries, entry)
	}

	return entries
}

func GetFeedByURL(url string) *Feed {
	rows, err := conn.Query("SELECT 1 FROM feeds WHERE url = ?", url)
	if err != nil {
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		// Feed exists
		return &Feed{URL: url}
	}
	return nil
}

// Retrieve all feeds
func GetFeeds() []*Feed {
	rows, err := conn.Query("SELECT id, url, title, description, last_updated FROM feeds")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var feeds []*Feed
	for rows.Next() {
		var id int
		var url, title, description string
		var lastUpdated string

		err := rows.Scan(&id, &url, &title, &description, &lastUpdated)
		if err != nil {
			return nil
		}

		feed := &Feed{
			ID:          int64(id),
			URL:         url,
			Title:       title,
			Description: description,
			LastUpdated: lastUpdated,
			Entries:     GetEntries(id),
		}
		feeds = append(feeds, feed)
		log.Println(len(feed.Entries), "entries loaded for feed:", feed.Title, " (", id, ")")
	}

	return feeds
}

// Mark an entry as read
func MarkRead(entryID int) error {
	stmt, err := conn.Prepare("UPDATE entries SET read = 1 WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(entryID)
	return err
}

// Update the last updated time for a feed
func MarkUpdated(feedID int64) error {
	stmt, err := conn.Prepare("UPDATE feeds SET last_updated = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(feedID)
	return err
}
