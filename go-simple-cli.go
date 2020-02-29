package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type response struct { // Structure for parsing response
	Main struct {
		Temp float32 `json:"temp"`
	} `json:"main"`
}

type weather struct {
	city string
	temp float32
}

func GetWeather(city string) weather {
	url := "http://api.openweathermap.org/data/2.5/weather"
	apikey := "83a404e36a0ffc57a902ebfd8f50480b"
	resp, err := http.Get(url + "?q=" + city + "&units=metric" + "&appid=" + apikey)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	obj := response{}

	if err := json.Unmarshal(body, &obj); err != nil {
		log.Fatal(err)
	}

	w := weather{city, obj.Main.Temp}

	return w
}
func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "city",
				Value: cli.NewStringSlice("Karachi", "Lahore"),
				Usage: "city to return temperature",
			},
		},
		Action: func(c *cli.Context) error {
			currWeather := make(chan weather)
			cities := c.StringSlice("city")

			for _, city := range cities {
				go func(c string) {
					w := GetWeather(c)
					currWeather <- w
				}(city)
			}
			for range cities {
				w := <-currWeather
				fmt.Printf("%s's temperature is %0.2f°C\n", w.city, w.temp)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
