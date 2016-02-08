# About xo #

xo is a cli tool to generate [Golang](https://golang.org/project/) types and
funcs based on a database schema.  xo is designed to vastly reduce the
overhead/redundancy in writing from scratch Go types and funcs for common
database tasks.

Note that while the code generated by xo is production quality, it is not the
goal, nor the intention for xo to be a "silver bullet" or to completely
eliminate writing SQL / Go code.

That said, xo can likely get you 99% "of the way there" on medium or large
database schemas and 100% of the way there for small or trivial database
schemas. In short, xo is a great launching point for developing standardized
packages for standard database abstractions/relationships.

Currently, xo only supports PostgreSQL, however there are plans to add
support for other database types, notably MySQL, Oracle, and SQLite.

## Design, Origin, and History ##

xo was originally developed while migrating a "large" application written in
PHP to Go. The schema in use in the original app, while well designed, had
become inconsistent over multiple iterations/generations, mainly due to
different naming styles adopted by various developers and database admins over
the preceding years.

This made for a relatively large headache in terms of unraveling old PHP code
as the code and API had drifted from the underlying names / fields as the
application code had also evolved through multiple rounds of different
development leads. Additionally, the code made use of multiple ORM-likes (most
notably Doctrine) and a custom, in-house semi-ORM-like code generator,
conceptually similar to this project.

As such, after a round of standardizing names, dropping accumulated cruft, and
adding a small number of relationship changes to the schema, the PHP code was
first fixed to match the schema changes. After that was determined to be a
success, the next target was a rewrite the backend services in Go.

In order to keep a similar and consistent work-flow for the developers, a code
generator similar to what was previously used with PHP was written for Go.
Additionally, at this time, but tangential to the story here, the API
definitions were ported from JSON to Protobuf to make use of the code
generation there.

xo is some of the fruits of those development efforts, and it is hoped that
others will be able to use and expand xo to support other databases (SQL or
otherwise).

Part of xo's goal is to avoid writing an ORM, or an ORM-like in Go, and to use
type-safe, fast, and idiomatic Go code. Additionally, the xo developers are of
the opinion that relational databases should have proper, well-designed
relationships and all the related definitions should reside within the database
schema itself. Call it a "self-documenting" schema. xo is an end to that
pursuit.

# Installation #

Install in the usual way for go:
```sh
go get -u github.com/knq/xo
```

# Usage #

Please note that xo is **NOT** an ORM. Rather, xo generates Go code by using
database metadata to query the types and relationships within the database, and
then generates representative Go types and funcs for well-defined database
relationships.

For example, given the following schema:
```PLpgSQL
CREATE TABLE authors (
  author_id SERIAL PRIMARY KEY,
  isbn text NOT NULL DEFAULT '' UNIQUE,
  name text NOT NULL DEFAULT '',
  subject text NOT NULL DEFAULT ''
);

CREATE INDEX authors_name_idx ON authors(name);

CREATE TABLE books (
    book_id SERIAL PRIMARY KEY,
    author_id integer NOT NULL REFERENCES authors(author_id),
    title text NOT NULL DEFAULT '',
    year integer NOT NULL DEFAULT 2000
);

CREATE INDEX books_title_idx ON books(title, year);
```

xo will generate the following (note: this is an abbreviated copy of actual
output. Please see the [example](example) folder for more a full example):
```go
// Author represents a row from public.authors.
type Author struct {
    AuthorID int    // author_id
    Isbn     string // isbn
    Name     string // name
    Subject  string // subject
}

// Exists determines if the Author exists in the database.
func (a *Author) Exists() bool { /* ... */ }

// Deleted provides information if the Author has been deleted from the database.
func (a *Author) Deleted() bool { /* ... */ }

// Insert inserts the Author to the database.
func (a *Author) Insert(db XODB) error { /* ... */ }

// Update updates the Author in the database.
func (a *Author) Update(db XODB) error { /* ... */ }

// Save saves the Author to the database.
func (a *Author) Save(db XODB) error { /* ... */ }

// Upsert performs an upsert for Author.
func (a *Author) Upsert(db XODB) error { /* ... */ }

// Delete deletes the Author from the database.
func (a *Author) Delete(db XODB) error { /* ... */ }

// AuthorByIsbn retrieves a row from public.authors as a Author.
//
// Looks up using index authors_isbn_key.
func AuthorByIsbn(db XODB, isbn string) (*Author, error) { /* ... */ }

// AuthorsByName retrieves rows from public.authors, each as a Author.
//
// Looks up using index authors_name_idx.
func AuthorsByName(db XODB, name string) ([]*Author, error) { /* ... */ }

// AuthorByAuthorID retrieves a row from public.authors as a Author.
//
// Looks up using index authors_pkey.
func AuthorByAuthorID(db XODB, authorID int) (*Author, error) { /* ... */ }

// Book represents a row from public.books.
type Book struct { /* ... */ }
    BookID   int    // book_id
    AuthorID int    // author_id
    Title    string // title
    Year     int    // year
}

// Exists determines if the Book exists in the database.
func (b *Book) Exists() bool { /* ... */ }

// Deleted provides information if the Book has been deleted from the database.
func (b *Book) Deleted() bool { /* ... */ }

// Insert inserts the Book to the database.
func (b *Book) Insert(db XODB) error { /* ... */ }

// Update updates the Book in the database.
func (b *Book) Update(db XODB) error { /* ... */ }

// Save saves the Book to the database.
func (b *Book) Save(db XODB) error { /* ... */ }

// Upsert performs an upsert for Book.
func (b *Book) Upsert(db XODB) error { /* ... */ }

// Delete deletes the Book from the database.
func (b *Book) Delete(db XODB) error { /* ... */ }

// Book returns the Author associated with the Book's AuthorID (author_id).
func (b *Book) Author(db XODB) (*Author, error) { /* ... */ }

// BookByBookID retrieves a row from public.books as a Book.
//
// Looks up using index books_pkey.
func BookByBookID(db XODB, bookID int) (*Book, error) { /* ... */ }

// BooksByTitle retrieves rows from public.books, each as a Book.
//
// Looks up using index books_title_idx.
func BooksByTitle(db XODB, title string, year int) ([]*Book, error) { /* ... */ }

// XODB is the common interface for database operations that can be used with
// types from public.
//
// This should work with database/sql.DB and database/sql.Tx.
type XODB interface {
    Exec(string, ...interface{}) (sql.Result, error)
    Query(string, ...interface{}) (*sql.Rows, error)
    QueryRow(string, ...interface{}) *sql.Row
}
```
