package link_manager

import (
	"errors"
	"github.com/the-gigi/delinkcious/pkg/db_util"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"time"

	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type DbLinkStore struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

const dbName = "link_manager"

func NewDbLinkStore(host string, port int, username string, password string) (store *DbLinkStore, err error) {
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
	store = &DbLinkStore{db, sb}
	return
}

func createSchema(db *sql.DB) (err error) {
	schema := `
        CREATE TABLE IF NOT EXISTS links (
          id SERIAL   PRIMARY KEY,
		  username    TEXT,
          url TEXT    UNIQUE NOT NULL,
          title TEXT  UNIQUE NOT NULL,
		  description TEXT,
		  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP          
        );
		CREATE UNIQUE INDEX IF NOT EXISTS links_username_idx ON links(username);


        CREATE TABLE IF NOT EXISTS tags (
          id SERIAL PRIMARY KEY,
          link_id   INTEGER REFERENCES links(id) ON DELETE CASCADE,			
          name      TEXT		  
        );
        CREATE UNIQUE INDEX IF NOT EXISTS tags_name_idx ON tags(name);
    `

	_, err = db.Exec(schema)
	return
}

func (s *DbLinkStore) GetLinks(request om.GetLinksRequest) (result om.GetLinksResult, err error) {
	q := s.sb.Select("*").From("links").Join("tags ON links.id = tags.link_id")
	q = q.Where(sq.Eq{"username": request.Username}).OrderBy("created_at")
	if request.StartToken != "" {
		var createdAt time.Time
		createdAt, err = time.Parse(time.RFC3339, request.StartToken)
		if err != nil {
			return
		}

		q = q.Where(sq.Gt{"created_at": createdAt})
	}
	if request.Tag != "" {
		q = q.Where(sq.Eq{"tag": request.Tag})
	}

	rows, err := q.RunWith(s.db).Query()
	if err != nil {
		return result, err
	}

	links := map[string]om.Link{}

	var link om.Link
	var id int
	var tag_id int
	var tag_name string
	var username string
	for rows.Next() {
		err = rows.Scan(&id, &username, &link.Url, &link.Title, &link.Description, &link.CreatedAt, &link.UpdatedAt, &tag_id, &id, &tag_name)
		if err != nil {
			return
		}

		_, ok := links[link.Url]
		if !ok {
			links[link.Url] = link
			result.Links = append(result.Links, link)
		}
	}

	result.NextPageToken = link.CreatedAt.UTC().Format(time.RFC3339)
	return
}

func (s *DbLinkStore) AddLink(request om.AddLinkRequest) (link *om.Link, err error) {
	link = &om.Link{
		Tags: map[string]bool{},
	}
	cmd := s.sb.Insert("links").Columns("username", "url", "title", "description").
		Values(request.Username, request.Url, request.Title, request.Description)
	_, err = cmd.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	q := s.sb.Select("*").From("links").Where(sq.Eq{"username": request.Username, "url": request.Url})
	var link_id int
	var username string
	row := q.RunWith(s.db).QueryRow()
	err = row.Scan(&link_id, &username, &link.Url, &link.Title, &link.Description, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		return
	}

	for t, _ := range request.Tags {
		cmd := s.sb.Insert("tags").Columns("link_id", "name").Values(link_id, t)
		_, err = cmd.RunWith(s.db).Exec()
		if err != nil {
			return
		}

		link.Tags[t] = true
	}

	return
}

func (s *DbLinkStore) UpdateLink(request om.UpdateLinkRequest) (link *om.Link, err error) {
	q := s.sb.Update("links").Where(sq.Eq{"username": request.Username, "url": request.Url})
	if request.Title != "" {
		q = q.Set("title", request.Title)
	}

	if request.Description != "" {
		q = q.Set("description", request.Description)
	}

	q = q.Suffix("RETURNING \"id\"")
	res, err := q.RunWith(s.db).Exec()
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected != 1 {
		err = errors.New("update failed")
	}

	var link_id int
	q.QueryRow().Scan(&link_id)

	for t, _ := range request.RemoveTags {
		_, err = s.sb.Delete("tags").Where(sq.Eq{"link_id": link_id, "name": t}).RunWith(s.db).Exec()
		if err != nil {
			return
		}

	}

	for t, _ := range request.AddTags {
		_, err = s.sb.Insert("tags").Columns("link_id", "name").Values(link_id, t).RunWith(s.db).Exec()
		if err != nil {
			return
		}
	}

	return

}

func (s *DbLinkStore) DeleteLink(username string, url string) (err error) {
	_, err = s.sb.Delete("links").Where(sq.Eq{"username": username, "url": url}).RunWith(s.db).Exec()
	return
}
