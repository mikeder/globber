package models

import "context"

// EntryArchive returns minimal information about all entries in the database.
func EntryArchive(ctx context.Context, db XODB) ([]*Entry, error) {
	const q string = `SELECT id, 
		author_id, slug, title, published, updated 
		FROM entries 
		ORDER BY published 
		DESC`

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	var entries []*Entry
	for rows.Next() {
		var e Entry

		err = rows.Scan(
			&e.ID,
			&e.AuthorID,
			&e.Slug,
			&e.Title,
			&e.Published,
			&e.Updated,
		)
		if err != nil {
			return entries, err
		}
		entries = append(entries, &e)
	}
	return entries, nil
}

// EntriesByPage retrieves up to 5 entries from the database with offset by page
// ordered by the published date.
func EntriesByPage(ctx context.Context, db XODB, page int) ([]*Entry, error) {
	var offset int
	if page > 0 {
		offset = page * 5
	}
	results, err := db.QueryContext(ctx, "SELECT * FROM entries ORDER BY published DESC LIMIT ?, 5", offset)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var entries []*Entry
	for results.Next() {
		var e Entry
		err = results.Scan(
			&e.ID,
			&e.AuthorID,
			&e.Slug,
			&e.Title,
			&e.Markdown,
			&e.HTML,
			&e.Published,
			&e.Updated,
		)
		if err != nil {
			return entries, err
		}
		entries = append(entries, &e)
	}
	return entries, nil
}
