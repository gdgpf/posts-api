package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Post struct {
	ID          int      `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	Description string   `json:"description" db:"description"`
	Latitude    *float32 `json:"latitude" db:"latitude"`
	Longitude   *float32 `json:"longitude" db:"longitude"`
	Likes       int      `json:"likes" db:"likes"`
	Dislikes    int      `json:"dislikes" db:"dislikes"`
}

func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {
		post := &Post{}
		postCreated := &Post{}

		if err := json.NewDecoder(c.Request.Body).Decode(post); err != nil {
			c.JSON(500, err)
			return
		}

		db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/dbdev")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.QueryRowx(
			"INSERT INTO post (title, description, latitude, longitude) VALUES ($1, $2, $3, $4) returning *",
			post.Title,
			post.Description,
			post.Latitude,
			post.Longitude,
		).StructScan(&postCreated)

		if err != nil {
			c.JSON(500, err)
			return
		}

		c.JSON(200, postCreated)
	})

	router.GET("/post/:id/like", func(c *gin.Context) {
		post := &Post{}
		id, _ := strconv.Atoi(c.Param("id"))

		db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/dbdev")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.QueryRowx(
			"UPDATE post SET likes = likes + 1 WHERE id = $1 returning *;",
			id,
		).StructScan(&post)

		if err != nil {
			c.JSON(500, err)
			return
		}

		c.JSON(200, post)
	})

	router.GET("/post/:id/dislike", func(c *gin.Context) {
		post := &Post{}
		id, _ := strconv.Atoi(c.Param("id"))

		db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/dbdev")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.QueryRowx(
			"UPDATE post SET dislikes = dislikes + 1 WHERE id = $1 returning *;",
			id,
		).StructScan(&post)

		if err != nil {
			c.JSON(500, err)
			return
		}

		c.JSON(200, post)
	})

	router.GET("/post", func(c *gin.Context) {
		posts := []Post{}

		db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/dbdev")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.Select(
			&posts,
			"SELECT * from post;",
		)

		if err != nil {
			c.JSON(500, err)
			return
		}

		c.JSON(200, posts)
	})

	router.Run(":10000")
}
