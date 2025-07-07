package feed

import (
	"database/sql"

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

func InitDB() {
	println("Initializing database...")
	// Initialize the SQLite database connection
	var err error
	conn, err = sql.Open("sqlite3", config.Config.DBPath)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		feed_id INTEGER,
		url TEXT,
		title TEXT,
		description TEXT,
		content TEXT,
		date_published DATETIME,
		read INTEGER DEFAULT 0,
		FOREIGN KEY(feed_id) REFERENCES feeds(id)
	)`)

	if err != nil {
		panic(err)
	}
}

func AddFeed(feed *gofeed.Feed) error {
	// Check if the feed already exists
	var exists bool
	var id int64
	var stmt *sql.Stmt
	var res sql.Result
	var err error

	err = conn.QueryRow("SELECT EXISTS(SELECT 1 FROM feeds WHERE url = ?)", feed.Link).Scan(&exists)
	if err != nil {
		println("Error checking if feed exists:", err.Error())
		return err
	}

	if !exists {
		println("Adding new feed to DB: ", feed.Link)
		// Insert  new feed into the database
		stmt, err = conn.Prepare("INSERT INTO feeds (url, title, description) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		res, err = stmt.Exec(feed.Link, feed.Title, feed.Description)

		if err != nil {
			println("Error inserting feed:", err)
			return err
		}
	}

	id, err = res.LastInsertId()
	for _, item := range feed.Items {
		err = AddEntry(id, item.Link, item.Title, item.Description, item.PublishedParsed.Format("2006-01-02 15:04:05"))
		if err != nil {
			println("Error adding entry:", err)
			return err
		}
	}
	return err
}

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

func GetEntries(feedID int) []*Entry {
	// Retrieve entries for a specific feed
	rows, err := conn.Query("SELECT id, url, title, description, date_published, read FROM entries WHERE feed_id = ?", feedID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var entries []*Entry
	for rows.Next() {
		var id int
		var url, title, description, content string
		var datePublished string
		var read int

		err := rows.Scan(&id, &url, &title, &description, &datePublished, &read, &content)
		if err != nil {
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

func GetFeeds() []*Feed {
	// Retrieve all feeds
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
	}

	return feeds
}

func MarkRead(entryID int) error {
	// Mark an entry as read
	stmt, err := conn.Prepare("UPDATE entries SET read = 1 WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(entryID)
	return err
}
