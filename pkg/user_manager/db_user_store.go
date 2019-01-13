package user_manager

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/the-gigi/delinkcious/pkg/db_util"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"math/rand"
	"strconv"
)

type DbUserStore struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

const dbName = "user_manager"

func NewDbUserStore(host string, port int, username string, password string) (store *DbUserStore, err error) {
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
	store = &DbUserStore{db, sb}
	return
}

func createSchema(db *sql.DB) (err error) {
	schema := `
        CREATE TABLE IF NOT EXISTS users (
          id SERIAL   PRIMARY KEY,
		  name    TEXT NOT NULL,
          email 	  TEXT UNIQUE NOT NULL
        );
		CREATE UNIQUE INDEX IF NOT EXISTS users_name_idx ON users(name);

        CREATE TABLE IF NOT EXISTS sessions (
          id SERIAL   PRIMARY KEY,
          user_id     INTEGER REFERENCES users(id) ON DELETE CASCADE,
		  session     TEXT NOT NULL,
          created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
		CREATE UNIQUE INDEX IF NOT EXISTS sessions_user_id_idx ON sessions(user_id);
		CREATE UNIQUE INDEX IF NOT EXISTS sessions_session_idx ON sessions(session);

    `

	_, err = db.Exec(schema)
	return
}

func (s *DbUserStore) Register(user om.User) (err error) {
	cmd := s.sb.Insert("users").Columns("name", "email").Values(user.Name, user.Email)
	_, err = cmd.RunWith(s.db).Exec()
	return
}

func (s *DbUserStore) Login(username string, authToken string) (session string, err error) {
	q := s.sb.Select("id").From("users").Where(sq.Eq{"name": username})
	_, err = q.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	var user_id int
	err = q.RunWith(s.db).QueryRow().Scan(&user_id)
	if err != nil {
		return
	}

	session = strconv.Itoa(rand.Int())
	cmd := s.sb.Insert("sessions").Columns("user_id", "session").Values(user_id, session)
	_, err = cmd.RunWith(s.db).Exec()
	return
}

func (s *DbUserStore) Logout(username string, session string) (err error) {
	q := s.sb.Select("id").From("users").Where(sq.Eq{"name": username})
	_, err = q.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	var user_id int
	q.QueryRow().Scan(&user_id)
	if err != nil {
		return
	}

	cmd := s.sb.Delete("sessions").Where(sq.Eq{"user_id": user_id, "session": session})
	_, err = cmd.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	cmd = s.sb.Delete("users").Where(sq.Eq{"id": user_id})
	_, err = cmd.RunWith(s.db).Exec()
	return
}
