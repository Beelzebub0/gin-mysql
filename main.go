package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/domesticated_pets")
	err = db.Ping()
	if err != nil {
		panic("Error")
	}
	defer db.Close()

	router := gin.Default()

	type Pets struct {
		Id_pets     int    `json: "id"`
		Name_pets   string `json: "name"`
		Gender_pets string `json: "gender"`
		Owner       string `json: "owner"`
	}

	// Menampilkan Detail Data Berdasarkan ID
	router.GET("/:id", func(c *gin.Context) {
		var (
			cats   Pets
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select * from cats where id = ?;", id)
		err = row.Scan(&cats.Id_pets, &cats.Name_pets, &cats.Gender_pets, &cats.Owner)
		if err != nil {
			// If no results send null
			result = gin.H{
				"Result": "No Data",
			}
		} else {
			result = gin.H{
				"Result": cats,
				"Total":  3,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET all persons
	router.GET("/", func(c *gin.Context) {
		var (
			cats  Pets
			catss []Pets
		)
		rows, err := db.Query("select * from cats;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&cats.Id_pets, &cats.Name_pets, &cats.Gender_pets, &cats.Owner)
			catss = append(catss, cats)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"Result": catss,
			"Total":  len(catss),
		})
	})

	// POST new person details
	router.POST("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		name := c.PostForm("name")
		gender := c.PostForm("gender")
		owner := c.PostForm("owner")
		stmt, err := db.Prepare("insert into cats (id, name, gender, owner) values(?,?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id, name, gender, owner)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(gender)
		buffer.WriteString(" ")
		buffer.WriteString(owner)
		buffer.WriteString(" ")
		defer stmt.Close()
		data := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Message": fmt.Sprintf(" Succesfully adding data %s ", data),
		})
	})

	// PUT - update a person details
	router.PUT("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		name := c.PostForm("name")
		gender := c.PostForm("gender")
		owner := c.PostForm("owner")
		stmt, err := db.Prepare("update cats set name= ?, gender = ?, owner= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(name, gender, owner, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(gender)
		buffer.WriteString(" ")
		buffer.WriteString(owner)
		buffer.WriteString(" ")
		defer stmt.Close()
		data := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Message": fmt.Sprintf("Succesfully changed id %s to %s", id, data),
		})
	})

	// Delete resources
	router.DELETE("/", func(c *gin.Context) {
		id := c.PostForm("id")
		stmt, err := db.Prepare("delete from hewan where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Message": fmt.Sprintf("Succes delet %s", id),
		})
	})
	router.Run(":8080")
}
