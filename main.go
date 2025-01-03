package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
)

var r *gin.Engine
var dbpool *pgxpool.Pool
var configEditor, colorsEditor *viper.Viper

func main() {
	configEditor = viper.New()
	configEditor.SetConfigName("config")
	configEditor.SetConfigType("yaml")
	configEditor.AddConfigPath("./config")
	configEditor.ReadInConfig()

	colorsEditor = viper.New()
	colorsEditor.SetConfigName("colors")
	colorsEditor.SetConfigType("json")
	colorsEditor.AddConfigPath("./config")
	colorsEditor.ReadInConfig()

	var err error
	dbpool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	authSetup()
	createTables()
	searchSetup()
	gob.Register(map[string]interface{}{})

	r = gin.Default()
	r.SetFuncMap(template.FuncMap{
		"contains": slices.Contains[[]string, string],
		"replace": func(toEscape string) string {
			return strings.ReplaceAll(toEscape, " ", "%20")
		},
	})
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/style.css", "./css/output.css")
	r.StaticFile("/icon.png", "./static/icon.png")
	r.LoadHTMLGlob("templates/*")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	loginRoutes()
	adminRoutes()
	voteRoutes()
	profileRoutes()

	r.GET("/", authMiddleware(), func(c *gin.Context) {
		session := sessions.Default(c)
		query := c.DefaultQuery("q", "")
		candidates := search(query)
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		flashes := session.Flashes()
		session.Save()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"text":    candidates,
			"flashes": flashes,
		})
	})
	r.Run()
}
