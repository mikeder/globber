package blog

import (
	"context"
	"database/sql"
	"time"

	"github.com/mikeder/globber/internal/models"
)

// Entry contains the fields for a blog entry.
type Entry struct {
	*models.Entry
}

// Store is where we keep blog posts and related information.
type Store struct {
	db *sql.DB
}

// New returns a pointer to a blog Store.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetEntriesByPage returns up to 5 Entries from the database, offset by page.
func (s *Store) GetEntriesByPage(ctx context.Context, page int) ([]Entry, error) {
	var posts []Entry

	entries, err := models.EntriesByPage(ctx, s.db, page)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		posts = append(posts, Entry{entry})
	}

	return posts, nil
}

// GetEntryBySlug returns a single Entry by its slug.
func (s *Store) GetEntryBySlug(ctx context.Context, slug string) (Entry, error) {
	entry, err := models.EntryBySlug(ctx, s.db, slug)
	if err != nil {
		return Entry{}, err
	}
	return Entry{entry}, nil
}

// GetArchive returns archived post data.
func (s *Store) GetArchive(ctx context.Context) ([]Entry, error) {
	archived, err := models.EntryArchive(ctx, s.db)
	if err != nil {
		return nil, err
	}

	var archive []Entry
	for _, entry := range archived {
		archive = append(archive, Entry{entry})
	}
	return archive, nil
}

// PostEntry adds a new Entry to the database.
func (s *Store) PostEntry(ctx context.Context, e *Entry) error {
	if e.Published.IsZero() {
		e.Published = time.Now()
	}
	e.Updated = time.Now()

	if err := e.Save(ctx, s.db); err != nil {
		return err
	}
	return nil
}

// GetEntryByID gets an Entry from the database by ID.
func (s *Store) GetEntryByID(ctx context.Context, id int) (*Entry, error) {
	entry, err := models.EntryByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	return &Entry{entry}, nil
}
