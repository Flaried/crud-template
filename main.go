package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// connection to the database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/users", getUsers(db))
	e.GET("/users/:id", getUser(db))
	e.POST("/users", createUser(db))
	e.PUT("/users", updateUser(db))
	e.DELETE("/users/:id", deleteUser(db))

	e.Logger.Fatal(e.Start(":1323"))
}

func jsonContentTypeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		return next(c)
	}
}

// get all users
func getUsers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := db.Query("SELECT id, name, email FROM users")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			users = append(users, user)
		}

		return c.JSON(http.StatusOK, users)
	}
}

// get user by id
func getUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get id from the URL
		id := c.Param("id")

		var user User
		row := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id)

		err := row.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, user)
	}
}

// create user from post
func createUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		_, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, user)
	}
}

// update user
func updateUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, user.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, user)
	}
}

// delete user
func deleteUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
