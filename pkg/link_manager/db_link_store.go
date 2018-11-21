package link_manager

import (
	"fmt"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"time"

	"database/sql"
	_ "github.com/lib/pq"
	sq "github.com/Masterminds/squirrel"
)

type DbLinkStore struct {
	db *sql.DB
}

func NewDbLinkStore(host string, port int, username string, password string) (store *DbLinkStore, err error) {
	mask := "host=%s port=%d user=%s password=%s dbname=link_manager sslmode=disable"
	dcn := fmt.Sprintf(mask, host, port, username, password)
	db, err := sql.Open("postgres", dcn)
	if err != nil {
		return
	}

	err = createSchema(db)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}
	store = &DbLinkStore{db}
	return
}

func createSchema(db *sql.DB) (err error) {
	schema := `
        CREATE TABLE IF NOT EXISTS links (
          id SERIAL PRIMARY KEY,
		  username TEXT,
          url TEXT UNIQUE NOT NULL,
          title TEXT UNIQUE NOT NULL,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
          description TEXT
        );
		CREATE UNIQUE INDEX IF NOT EXISTS username_idx ON links(username);


        CREATE TABLE IF NOT EXISTS tags (
          id SERIAL PRIMARY KEY,
          link_id   INTEGER FOREIGN KEY links(id)			
          name  TEXT,		  
        );
        CREATE UNIQUE INDEX IF NOT EXISTS tag_link_name_idx ON tags(link_id, name);
    `

	_, err = db.Exec(schema)
	return
}

func (s *DbLinkStore) GetLinks(request om.GetLinksRequest) (result om.GetLinksResult, err error) {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := sb.Select("*").From("links").Join("tags USING link_id").OrderBy("created_at")
	if request.StartToken != "" {
		createdAt, err := time.Parse(time.RFC3339, request.StartToken)
		if err != nil {
			return
		}

		q = q.Where(sq.Gt{"created_at": createdAt})
	}
	if request.Tag != "" {
		q = q.Where(sq.Eq{"tag": request.Tag})
	}

	fmt.Println(q.ToSql())
	rows, err := q.RunWith(s.db).Query()

	links := map[string]om.Link{}

	var link om.Link
	var id int
	for rows.Next() {
		err = rows.Scan(&id, &link.Url, &link.Title, &link.Description, &link.CreatedAt, &link.UpdatedAt)
		if err != nil {
			return
		}

		_, ok := links[link.Url]
		if !ok {
			links[link.Url] = link
		}
	}

	result.NextPageToken = link.CreatedAt.UTC().Format(time.RFC3339)
	return
}

func (s *DbLinkStore) AddLink(request om.AddLinkRequest) (link *om.Link, err error) {
	cmd := sq.Insert("links").Columns("username", "url", "title", "description").
		                    Values(request.Username, request.Url, request.Title, request.Description)
	_, err = cmd.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := sb.Select("*").From("links").Where(sq.Eq{"url": request.Url})
	fmt.Println(q.ToSql())
	var id int
	err := q.RunWith(s.db).QueryRow().Scan(&id, &link.Url, &link.Title, &link.Description, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		return
	}

	for _, t := range request.Tags {
		??????
	}



	return
}

func (s *DbLinkStore) UpdateLink(request om.UpdateLinkRequest) (link *om.Link, err error) {
	return
}

func (s *DbLinkStore) DeleteLink(username string, url string) error {
	return nil
}
