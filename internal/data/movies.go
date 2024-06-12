package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/TomiwaAribisala-git/greenlight.git/internal/validator"
	"github.com/lib/pq"
)

// fields in capital letter for them to be exported(visible) to Go encoding/json package
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty,string"` // advanced json customization
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version,omitempty,string"` // the string directive only work on struct fields which have int*/uint*/float*/bool types
}

// the fields and types in the movie struct above is analogous to the fields and types of the database migrations table

// Define a MovieModel struct type which wraps a sql.DB connection pool
type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
	INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	// You can also use the pq.Array() adapter function in the same way with []bool,
	// []byte, []int32, []int64, []float32 and []float64 slices
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() and pass the context as the first argument.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	// smallint, smallserial int16
	// integer, serial int32
	// bigint, bigserial int64

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`

	var movie Movie

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	// Implementing optimistic locking, incrementing the version value as part of the query
	// Using other locking field type: SET title = $1, year = $2, runtime = $3, genres = $4, version = uuid_generate_v4()
	query := `
	UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1		
	WHERE id = $5 AND version = $6	
	RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict // record version-wise has changed or deleted
		default:
			return err
		}
	}
	return nil
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM movies
	WHERE id = $1`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) { //([]*Movie, Metadata, error)
	/*
		// Construct the SQL query to retrieve all movie records
		query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		ORDER BY id`

		// Filtering Data: for query string parameters
		// The @> symbol is the ‘contains’ operator for PostgreSQL arrays
		query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (LOWER(title) = LOWER($1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY id`

		// Sorting Lists
		// Add an ORDER BY clause and interpolate the sort column and direction. Importantly
		// notice that we also include a secondary sort on the movie ID to ensure a
		// consistent ordering.
		query := fmt.Sprintf(`
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())

		// Paginating Lists
		// Update the SQL query to include the LIMIT and OFFSET clauses with placeholder
		// parameter values.
		query := fmt.Sprintf(`
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

		// As our SQL query now has quite a few placeholder parameters, let's collect the
		// values for the placeholders in a slice. Notice here how we call the limit() and
		// offset() methods on the Filters struct to get the appropriate values for the
		// LIMIT and OFFSET clauses.
		args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}

		// And then pass the args slice to QueryContext() as a variadic parameter.
		rows, err := m.DB.QueryContext(ctx, query, args...)

		// Pagination Metadata
		query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

		args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
		rows, err := m.DB.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}

		// Declare a totalRecords variable.
		totalRecords := 0
		movies := []*Movie{}

		for rows.Next() {
			var movie Movie

			err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
			)
			if err != nil {
				return nil, Metadata{}, err // Update this to return an empty Metadata struct.
			}

			movies = append(movies, &movie)
		}

		if err = rows.Err(); err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}

		// Generate a Metadata struct, passing in the total record count and pagination
		// parameters from the client.
		metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

		// Include the metadata struct when returning.
		return movies, metadata, nil
		}
	*/

	// PostgreSQL’s full textsearch functionality, to use by adapting it to
	// support partial matches, rather than requiring a match on the full title.
	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (genres @> $2 OR $2 = '{}')
	ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// rows, err := m.DB.QueryContext(ctx, query)
	// rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres))
	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []*Movie{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var movie Movie

		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}

		// Add the Movie struct to the slice.
		movies = append(movies, &movie)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK, then return the slice of movies.
	return movies, nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

/*
func (m MovieModel) GetAll(title string, year int32, runtime Runtime, genres []string) ([]*Movie, error) {

	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	ORDER BY id`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []*Movie{}

	for rows.Next() {

		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the slice of movies.
	return movies, nil
}
*/
