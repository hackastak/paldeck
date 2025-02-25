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
	Type string
	Description string
	Aura string
	AuraDescription string
	Suitability []int
}

// func newPal(ID int, Name, Key, Type, Description, Aura, AuraDescription string, Suitability []int) Pal {
// 	return Pal{
// 		ID: ID,
// 		Name: Name,
// 		Key: Key,
// 		Type: Type,
// 		Description: Description,
// 		Aura: Aura,
// 		AuraDescription: AuraDescription,
// 		Suitability: Suitability,
// 	}
// }

type Pals = []Pal

type PalData struct {
	Pals Pals
}

func queryPals(db *sql.DB) (Pals){
	rows, err := db.Query("SELECT * FROM pals")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	Pals := []Pal{}

	for rows.Next() {
		var pal Pal

		if err := rows.Scan(&pal.ID, &pal.Name, &pal.Key); err != nil {
			fmt.Println("Error scanning row: ", err)
			return nil
		}

		Pals = append(Pals, pal)
		fmt.Println("Pals Slice WIP: ", Pals)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration: ", err)
		return nil
	}

	return Pals
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

	data := queryPals(db)
	PalData := PalData{Pals: data}
	fmt.Println("Pals After Query: ", PalData)


	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()
	e.Static("/styles", "styles")
	e.Static("/assets", "assets")
 
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", PalData)
	})

	e.GET("/pals/new", func(c echo.Context) error {
		return c.Render(200, "createPal", PalData)
	})

	// e.POST("/pals/new", func(c echo.Context) error {
	// 	return c.JSON(200, palData)
	// })
	
	e.Logger.Fatal(e.Start(":8000"))

}