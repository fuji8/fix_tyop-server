package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/genproto/googleapis/type/dayofweek"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var tables = []interface{}{
	Dialogue{},
}

const aitalk = "https://webapi.aitalk.jp/webapi/v5/ttsget.php"

type Dialogue struct {
	ID        uuid.UUID `gorm:"type:char(36); primaryKey"`
	DayOfWeek int32
	Oder      int
	Dialogue  string
}

func (d *Dialogue) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID != uuid.Nil {
		return nil
	}
	d.ID, err = uuid.NewV4()
	return err
}

type DialogueReq struct {
	Dialogue string `json:"dialogue"`
}

var CharacterName = map[dayofweek.DayOfWeek]string{
	dayofweek.DayOfWeek_MONDAY:    "aoi_emo",
	dayofweek.DayOfWeek_TUESDAY:   "akane_west_emo",
	dayofweek.DayOfWeek_WEDNESDAY: "yuzuru_emo",
}

func main() {

	password := os.Getenv("mariadb_password")
	if password == "" {
		password = "password"
	}
	host := os.Getenv("")
	if host == "" {
		host = "localhost"
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true&loc=Local", "root", password, host, "app"),
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   false, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(tables...)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Golang, Hello, World!")
	})

	e.GET("/audios", func(c echo.Context) error {
		dow := c.QueryParam("dayOfWeek")
		dayOfWeek, ok := dayofweek.DayOfWeek_value[dow]
		if !ok {
			fmt.Println("###: ", dow)
			return echo.NewHTTPError(http.StatusNotFound)
		}
		order, err := strconv.Atoi(c.QueryParam("order"))
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		u, _ := url.Parse(aitalk)
		q := setAuth(u)
		q.Set("ext", "mp3")
		q.Set("speaker_name", CharacterName[dayofweek.DayOfWeek(dayOfWeek)])

		var dialogue Dialogue
		err = db.Where("day_of_week = ?", dayOfWeek).Where("oder = ?", order).Take(&dialogue).Error
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		q.Set("text", dialogue.Dialogue)
		u.RawQuery = q.Encode()

		resp, err := http.Get(u.String())
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		io.Copy(buf, resp.Body)

		return c.Blob(http.StatusOK, "audio/mpeg", buf.Bytes())
	})

	e.POST("/dialogues", func(c echo.Context) error {
		dow := c.QueryParam("dayOfWeek")
		dayOfWeek, ok := dayofweek.DayOfWeek_value[dow]
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		order, err := strconv.Atoi(c.QueryParam("order"))
		if err != nil || order != 5 {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		var req DialogueReq
		if err = c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		err = db.Create(&Dialogue{
			DayOfWeek: dayOfWeek,
			Oder:      order,
			Dialogue:  req.Dialogue,
		}).Error

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func setAuth(u *url.URL) url.Values {
	q := u.Query()
	q.Set("username", os.Getenv("aitalk_username"))
	q.Set("password", os.Getenv("aitalk_password"))
	return q
}
