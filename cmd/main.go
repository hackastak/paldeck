package main

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

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
	Key int
	Name string
	ImgUrl string
}



func main () {
	fmt.Println("Hey! It's fuck you buddy!")
	tursoUrl := "libsql://paldeck-hackastak.turso.io?authToken=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MTA3MzE0NDcsImlkIjoiOWZiYmVjZWItMWI2ZC00Yzg0LTgzNmMtYjRlYmZlMmQxMTc4In0.pIJD1pTJYGvJX9sb_Z3iuwqKJuJrYwFuldj32k1uIRMwjUlmd1-lInXmjiY8oGJWzO5RUWYXvtcANgGd14y1BA"

	db, err := sql.Open("libsql", tursoUrl)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to open db %s: %s", tursoUrl, err)
    os.Exit(1)
  }
  defer db.Close()


	e := echo.New()
	e.Use(middleware.Logger())

	pal := Pal {Name: "Rubio"}
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", pal)
	})
	
	e.Logger.Fatal(e.Start(":8000"))

}