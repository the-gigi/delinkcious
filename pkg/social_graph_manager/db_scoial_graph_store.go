package social_graph_manager

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/the-gigi/delinkcious/pkg/db_util"
)

type DbSocialGraphStore struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

const dbName = "social_graph_manager"

func NewDbSocialGraphStore(host string, port int, username string, password string) (store *DbSocialGraphStore, err error) {
	db, err := db_util.EnsureDB(host, port, username, password, dbName)
	if err != nil {
		return
	}

	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	err = createSchema(db)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}
	store = &DbSocialGraphStore{db, sb}
	return
}

func createSchema(db *sql.DB) (err error) {
	schema := `
        CREATE TABLE IF NOT EXISTS social_graph (
          id SERIAL   PRIMARY KEY,
		  followed    TEXT NOT NULL,
          follower 	  TEXT NOT NULL,
		  UNIQUE (followed, follower)
        );
		CREATE INDEX IF NOT EXISTS social_graph_follower_idx ON social_graph(follower);
		CREATE INDEX IF NOT EXISTS social_graph_followed_idx ON social_graph(followed);
    `

	_, err = db.Exec(schema)
	return
}

func (s *DbSocialGraphStore) Follow(followed string, follower string) (err error) {
	cmd := s.sb.Insert("social_graph").Columns("followed", "follower").Values(followed, follower)
	_, err = cmd.RunWith(s.db).Exec()
	return
}

func (s *DbSocialGraphStore) Unfollow(followed string, follower string) (err error) {
	cmd := s.sb.Delete("social_graph").Where(sq.Eq{"followed": followed, "follower": follower})
	r, err := cmd.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected != 1 {
		return errors.New("unable to unfollow")
	}

	return
}

func (s *DbSocialGraphStore) GetFollowers(username string) (followers map[string]bool, err error) {
	followers = map[string]bool{}
	q := s.sb.Select("follower").From("social_graph").Where(sq.Eq{"followed": username})
	rows, err := q.RunWith(s.db).Query()
	if err != nil {
		return
	}

	follower := ""
	for rows.Next() {
		err = rows.Scan(&follower)
		if err != nil {
			return
		}

		followers[follower] = true
	}

	return
}

func (s *DbSocialGraphStore) GetFollowing(username string) (following map[string]bool, err error) {
	following = map[string]bool{}
	q := s.sb.Select("followed").From("social_graph").Where(sq.Eq{"follower": username})
	rows, err := q.RunWith(s.db).Query()
	if err != nil {
		return
	}

	followed := ""
	for rows.Next() {
		err = rows.Scan(&followed)
		if err != nil {
			return
		}

		following[followed] = true
	}

	return
}
