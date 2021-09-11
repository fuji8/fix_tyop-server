package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func saveAndServeSound(c echo.Context) error {
	text := c.QueryParam("text")
	log.Println(text)

	if err := getMusic(text); err != nil {
		formattedErr := fmt.Errorf("Error hello L:28 %w\n", err)
		log.Println(formattedErr.Error())
		return err
	}

	return c.File("hoge.mp3")
}

func getMusic(text string) error {
	url := fmt.Sprintf("だみー")

	out, err := os.Create("hoge.mp3")
	if err != nil {
		formattedErr := fmt.Errorf("Error getMusic L31: %w\n", err)
		log.Println(formattedErr.Error())
		return formattedErr
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		formattedErr := fmt.Errorf("Error getMusic L40: %w\n", err)
		log.Println(formattedErr.Error())
		return formattedErr
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		_, err := io.Copy(out, resp.Body)
		if err != nil {
			formattedErr := fmt.Errorf("Error getMusic L51: %w\n", err)
			log.Println(formattedErr.Error())
			return formattedErr
		}
	} else {
		formattedErr := fmt.Errorf("Error getMusic requests fatal: %d\n", resp.StatusCode)
		log.Println(formattedErr.Error())
		return formattedErr

	}

	return nil
}
