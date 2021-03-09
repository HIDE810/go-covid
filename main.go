package main

import (
	"context"
	"github.com/dustin/go-humanize"
	"github.com/itsksaurabh/go-corona"
	"github.com/labstack/echo"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}

	e := echo.New()
	e.Static("/", "")
	e.Renderer = t
	e.GET("/", index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Start(":" + port)
}

func index(c echo.Context) error {

	client := gocorona.Client{}
	ctx := context.Background()

	latest, err := client.GetLatestData(ctx)
	if err != nil {
		log.Fatal("request failed:", err)
	}

	info := map[string]interface{}{
		"confirmed": humanize.Comma(int64(latest.Data.Confirmed)),
		"deaths":    humanize.Comma(int64(latest.Data.Deaths)),
		"recovered": humanize.Comma(int64(latest.Data.Recovered)),
	}

	return c.Render(http.StatusOK, "index.html", info)
}
