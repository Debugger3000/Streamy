package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Series struct {
	Title    string `json:"Title"`
	Released string `json:"Released"` // first release
	Year     string `json:"Year"`     // total run years
	Rated    string `json:"Rated"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Poster   string `json:"Poster"`
}

// pre much exact match on title...
// one object

func GetSeriesInfo(seriesName string) Series {
	var api_key = os.Getenv("OMDB_API_KEY")

	fmt.Println(api_key)
	// http://www.omdbapi.com/?apikey=[yourkey]&

	url := fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&t=%s&type=series", api_key, seriesName)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// strings are treated as BYTES
	// if we index a string, we get their byte count

	var s Series
	if err := json.Unmarshal(body, &s); err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
	return s

}
