package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/mrizkyanditama/lelangin/api/middlewares"
	"gopkg.in/olahol/melody.v1"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/mrizkyanditama/lelangin/api/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
	Mu     *sync.Mutex
	Mel    *melody.Melody
}

var errList = make(map[string]string)

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	// If you are using mysql, i added support for you here(dont forgot to edit the .env file)
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	//database migration
	server.DB.AutoMigrate(
		&models.User{},
		&models.ResetPassword{},
		&models.Category{},
		&models.Product{},
		&models.Tag{},
		&models.Auction{},
		&models.Bid{},
	)

	// if err_migrate != nil {
	// 	log.Fatal("This is the error connecting to postgres:", err_migrate)
	// }

	server.Router = gin.Default()
	server.Router.Use(middlewares.CORSMiddleware())

	server.Mu = new(sync.Mutex)
	server.Mel = melody.New()

	server.initializeBidsController()
	server.initializeRoutes()

}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
