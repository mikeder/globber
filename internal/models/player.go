package models

import "context"

// AllPlayers returns all players from the database.
func AllPlayers(ctx context.Context, db XODB) ([]*Player, error) {
	const q string = `SELECT ` +
		`id, name, uuid, first_seen, last_seen ` +
		`FROM blog.players`

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	var players []*Player
	for rows.Next() {
		var p Player

		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.UUID,
			&p.FirstSeen,
			&p.LastSeen,
		)
		if err != nil {
			return players, err
		}
		players = append(players, &p)
	}
	return players, nil
}
