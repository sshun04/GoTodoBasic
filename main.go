package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

var dbDeviceName string = "sqlite3"
var dbFileName string = "test.sqlite3"

const key_text  = "text"
const key_status  =  "status"
const key_id  = "id"
const redirect_code = 302
const html_code  = 200

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	router.GET("/", func(context *gin.Context) {
		todos := dbGetAll()
		context.HTML(html_code, "index.html", gin.H{"todos": todos,})
	})

	router.POST("/new", func(context *gin.Context) {
		text := context.PostForm(key_text)
		status := context.PostForm(key_status)
		dbInsert(text,status)
		context.Redirect(redirect_code,"/")
	})

	router.GET("/detail/:id", func(context *gin.Context) {
		n := context.Param(key_id)
		id,err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		context.HTML(html_code,"detail.html",gin.H{"todo":todo})
	})

	router.POST("/update/:id", func(context *gin.Context) {
		n := context.Param(key_id)
		id,err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		text := context.PostForm(key_text)
		status := context.PostForm(key_status)
		dbUpdate(id,text,status)
		context.Redirect(redirect_code,"/")
	})

	router.GET("/delete_check/:id", func(context *gin.Context) {
		n := context.Param(key_id)
		id,err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		context.HTML(html_code,"delete.html",gin.H{"todo":todo})
	})

	router.POST("/delete/:id", func(context *gin.Context) {
		n := context.Param(key_id)
		id,err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbDelete(id)
		context.Redirect(redirect_code,"/")
	})

	router.Run()
}

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

func dbInit() {
	db, err := gorm.Open(dbDeviceName, dbFileName)
	if err != nil {
		panic("failure open database:Init")
	}
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

func dbInsert(text string, status string) {
	db, err := gorm.Open(dbDeviceName, dbFileName)
	if err != nil {
		panic("failure open database:Insert")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

func dbGetAll() []Todo {
	db, err := gorm.Open(dbDeviceName, dbFileName)
	if err != nil {
		panic("failure open database:GetAll")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

func dbGetOne(id int) Todo {
	db, err := gorm.Open(dbDeviceName, dbFileName)
	if err != nil {
		panic("failure open database:GetOne")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open(dbDeviceName, dbFileName)
	if err != nil {
		panic("failure open database:Update")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

func dbDelete(id int) {
	db, err := gorm.Open(dbDeviceName, dbFileName)
	if err != nil {
		panic("failure open database:Dlete")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}