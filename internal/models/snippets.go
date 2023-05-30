package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// SQL statements
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Use the Exec() method on the embedded connection pool to execute the statement.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}
	// Use the LastInsertId() method on the result on get the ID of newly inserted record in snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	// The ID returned has the type int64, convert it to an int type.
	return int(id), nil
}

// return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// QueryRow() method returns a pointer to a sql.Row object.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// Use errors.Is() function check error and return own ErrNoRecord error instead
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute SQL statement.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// rows.Close() to ensure the sql.Rows resultest is always properly closed
	// before the Latest() method returns.
	defer rows.Close()

	// Initlizzlize an empty slice to hold the Snippet structs.
	snippets := []*Snippet{}

	// Use rows.Next() to iterate through the rows in the resultset.
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	// rows.Next() loop has finished, call rows.Err() to retrieve any error that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
