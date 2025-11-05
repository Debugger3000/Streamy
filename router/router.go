package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type MetaData struct {
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

func GetSeriesInfo(seriesName string) MetaData {
	var api_key = os.Getenv("OMDB_API_KEY")

	// rturn empty meta data is api key cannot be grabbed
	if api_key == "" {
		log.Println("OMDB_API_KEY not set")
		return MetaData{}
	}

	fmt.Println(api_key)
	// http://www.omdbapi.com/?apikey=[yourkey]&
	encodedName := url.QueryEscape(seriesName)

	url := fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&t=%s&type=series", api_key, encodedName)

	resp, err := http.Get(url)
	// api call error, give empty meta data
	if err != nil {
		log.Println("OMDB request error:", err)
		return MetaData{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("OMDB read error:", err)
		return MetaData{}
	}

	// strings are treated as BYTES
	// if we index a string, we get their byte count

	var s MetaData
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println("OMDB parse error:", err)
		return MetaData{}
	}

	fmt.Println(s)
	return s
}

func GetMovieInfo(seriesName string) MetaData {
	var api_key = os.Getenv("OMDB_API_KEY")

	// rturn empty meta data is api key cannot be grabbed
	if api_key == "" {
		log.Println("OMDB_API_KEY not set")
		return MetaData{}
	}

	fmt.Println(api_key)
	// http://www.omdbapi.com/?apikey=[yourkey]&
	encodedName := url.QueryEscape(seriesName)

	url := fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&t=%s&type=series", api_key, encodedName)

	resp, err := http.Get(url)
	// api call error, give empty meta data
	if err != nil {
		log.Println("OMDB request error:", err)
		return MetaData{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("OMDB read error:", err)
		return MetaData{}
	}

	// strings are treated as BYTES
	// if we index a string, we get their byte count

	var s MetaData
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println("OMDB parse error:", err)
		return MetaData{}
	}

	fmt.Println(s)
	return s
}
