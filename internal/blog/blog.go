package blog

import (
	"context"
	"database/sql"
)

// Store is where we kep blog posts and related information.
type Store struct {
	db *sql.DB
}

// New returns a pointer to a blog Store.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// Post contains the fields for a blog post.
type Post struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Markdown  string `json:"markdown"`
	HTML      string `json:"html"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
}

// GetPostBySlug returns a *Post by its slug.
func (s *Store) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	result, err := s.db.QueryContext(ctx, "SELECT * FROM entries WHERE slug=?", slug)
	if err != nil {
		return nil, err
	}
	var post Post
	for result.Next() {
		err = result.Scan(
			&post.ID,
			&post.Author,
			&post.Slug,
			&post.Title,
			&post.Markdown,
			&post.HTML,
			&post.Published,
			&post.Updated,
		)
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPosts returns a []Post of all the posts from the database.
func (s *Store) GetPosts(ctx context.Context, page int) ([]Post, error) {
	var offset int
	if page > 0 {
		offset = page * 5
	}
	results, err := s.db.QueryContext(ctx, "SELECT * FROM entries ORDER BY published DESC LIMIT ?, 5", offset)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var posts []Post
	for results.Next() {
		var post Post
		err = results.Scan(
			&post.ID,
			&post.Author,
			&post.Slug,
			&post.Title,
			&post.Markdown,
			&post.HTML,
			&post.Published,
			&post.Updated,
		)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
