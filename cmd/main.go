package main

import (
	"fmt"
	"html/template"
	"database/sql"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_"github.com/tursodatabase/libsql-client-go/libsql"

)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Pal struct {
	ID int
	Name string
	Key string
}

var pals []Pal

func queryPals(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM pals")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var pal Pal

		if err := rows.Scan(&pal.ID, &pal.Name, &pal.Key); err != nil {
			fmt.Println("Error scanning row: ", err)
			return
		}

		pals = append(pals, pal)
		fmt.Printf("%d: %s -- %s\n", pal.ID, pal.Name, pal.Key)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration: ", err)
	}

}


func main () {

	os.Setenv("PALDECK_TURSO_AUTH_URL", "libsql://paldeck-hackastak.turso.io")
	os.Setenv("PALDECK_TURSO_AUTH_TOKEN", "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MTExOTU2NzksImlkIjoiOWZiYmVjZWItMWI2ZC00Yzg0LTgzNmMtYjRlYmZlMmQxMTc4In0._3Xa4Rc1fBiOgBTH6JGAZ8QbO0_z-szhc6rm2v0zW9HaPRpvcOl052X4sdPlEI8MZOIBOmlpB95_KiLud-U3CA")
	dbUrl := os.Getenv("PALDECK_TURSO_AUTH_URL")
	dbAuthToken := os.Getenv("PALDECK_TURSO_AUTH_TOKEN")

	// Turso Database connection configuration
	tursoUrl := dbUrl+"?authToken="+dbAuthToken
	db, err := sql.Open("libsql", tursoUrl)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to open db %s: %s", tursoUrl, err)
    os.Exit(1)
  }
  defer db.Close()

	queryPals(db)
	


	e := echo.New()
	e.Use(middleware.Logger())

	pal := pals[0]
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", pal)
	})
	
	e.Logger.Fatal(e.Start(":8000"))

}