package blog

import "database/sql"

type Store struct {
	db *sql.DB
}

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

func (s *Store) GetPost(id int) (*Post, error) {
	result, err := s.db.Query("SELECT * FROM entries WHERE id=?", id)
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

func (s *Store) GetPosts() ([]Post, error) {
	results, err := s.db.Query("SELECT * FROM entries ORDER BY published DESC")
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
