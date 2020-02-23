package blog

import (
	"context"
	"database/sql"

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

// GetEntriesByPage returns a []Post of all the posts from the database.
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

// GetEntryBySlug returns an *Entry by its slug.
func (s *Store) GetEntryBySlug(ctx context.Context, slug string) (Entry, error) {
	entry, err := models.EntryBySlug(ctx, s.db, slug)
	if err != nil {
		return Entry{}, err
	}
	return Entry{entry}, nil
}

// GetArchive returns a []Post of archived post data.
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
