package service

import (
	"github.com/gin-gonic/gin"
	"github.com/uber/h3-go/v4"
	"log"
	"net/http"
	"strconv"
	"strings"
	"training/internal/pkg/config"
	"training/internal/pkg/database"
)

const (
	configYamlFile string = "config.yaml"
)

type Arg struct {
	Conf *config.Base
	DB   *database.Database
}

type query struct {
	Search   string `json:"search"`
	Location string `json:"location"`
}

type row struct {
	H3Key string   `json:"h3index"`
	Word  []string `json:"words"`
}

func LoadConfig() *Arg {
	appConfig := config.Base{}
	err := appConfig.Load(configYamlFile)
	if err != nil {
		log.Print("config was empty!", err)
	}
	db, err := database.NewDatabase(&appConfig)
	if err != nil {
		log.Print("failed to establish a database", err)
	}
	arg := &Arg{&appConfig, db}
	return arg
}

func Engine(arg *Arg) *gin.Engine {
	router := gin.Default()
	router.GET(arg.Conf.GetAPI, arg.GetAlbums)
	router.POST(arg.Conf.PostAPI, arg.PostAlbums)
	return router
}

func (argument *Arg) GetAlbums(c *gin.Context) {
	db := argument.DB
	key := c.Param("id")
	val := db.Get(key)
	row := row{H3Key: key, Word: val}
	c.IndentedJSON(http.StatusOK, row)
}

func (argument *Arg) PostAlbums(c *gin.Context) {
	db := argument.DB
	conf := argument.Conf
	var newQuery query

	if err := c.BindJSON(&newQuery); err != nil {
		return
	}

	loc := strings.Split(newQuery.Location, ",")
	lat, err := strconv.ParseFloat(loc[0], 64)
	if err != nil {
		log.Print("converting Lat to float failed!", err)
	}
	lng, err := strconv.ParseFloat(loc[1], 64)
	if err != nil {
		log.Print("converting Lng to float failed!", err)
	}
	latLng := h3.NewLatLng(lat, lng)
	cell := h3.LatLngToCell(latLng, conf.Resolution)

	words := strings.Split(newQuery.Search, " ")
	err = db.Add(cell.String(), words)
	if err != nil {
		log.Print("error while adding to db", err)
	}
	row := row{H3Key: cell.String(), Word: db.Get(cell.String())}
	c.IndentedJSON(http.StatusCreated, row)
}
