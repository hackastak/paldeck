package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e := echo.New()
	e.Use(middleware.Logger())

	pal := Pal {Name: "Rubio"}
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", pal)
	})
	
	e.Logger.Fatal(e.Start(":8000"))

}