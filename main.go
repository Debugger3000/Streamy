package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"net/url"
	"streamy/router"

	"github.com/joho/godotenv"
)

// ffmpeg -i "C:\path\to\video.mkv" -c copy "C:\path\to\video.mp4"
// to convert a .mkv into a .mp4 for html video compatibility

// encoding for mobile friendly + desktop
// ffmpeg -i input.mp4 -c:v libx264 -pix_fmt yuv420p -crf 23 -preset fast -c:a aac -b:a 128k output.mp4
// h264_nvenc
// in place of 'libx264', since this one uses CPU and the other uses GPU

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is the Streamy server!")
}

// Helper function to extract season number from folder name
func extractSeasonNumber(name string) int {
	// Example: "Season 1" -> 1
	parts := strings.Fields(name)
	if len(parts) >= 2 {
		if num, err := strconv.Atoi(parts[1]); err == nil {
			return num
		}
	}
	return 0
}

// structs
// -----------------------------------------

// define type to hold video names...
// main page where lists of basic strings for titles are used...
type PageData struct {
	Movies  []string
	TVShows []string
	Other   []string
}

// ------------------------

// current directory media fileName to stream from that filename in /media
type MoviePageData struct {
	Filename  string
	Extension string
	router.MetaData
}

// ------
// TV show structs
type SeasonStruct struct {
	Season   int
	Episodes []Episode
}

type Episode struct {
	FileName  string
	Extension string
}

type ShowPageData struct {
	Content []SeasonStruct
	router.MetaData
}

// ------------------------------------------

// --------------------------------

// main function
func main() {

	// Load .env file so env vars become available to os.Getenv
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, continuing...")
	}

	// variables
	var port = os.Getenv("PORT")

	fmt.Println("Running main")

	http.HandleFunc("/test", handler)
	//

	// Serve media files
	fs := http.FileServer(http.Dir("./media"))              // serve from this directory
	http.Handle("/media/", http.StripPrefix("/media/", fs)) // serve that directory on path "/media/"

	// Serve static files (CSS, JS, etc.)
	fsStatic := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))

	tmpl := template.Must(template.ParseFiles("./templates/index.html"))      // serve index.html
	moviePage := template.Must(template.ParseFiles("./templates/movie.html")) // // serve movie html template
	showPage := template.Must(template.ParseFiles("./templates/show.html"))

	// serve video js file
	// Serve the video-js folder at /video-js-8.23.6
	video_js := http.FileServer(http.Dir("./video-js-8.23.6"))
	http.Handle("/video-js-8.23.6/", http.StripPrefix("/video-js-8.23.6/", video_js))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Read files from media folder
		movieEntries, err := os.ReadDir("./media/movies")
		if err != nil {
			http.Error(w, "Unable to read media folder", http.StatusInternalServerError)
			return
		}

		var movies []string
		for _, entry := range movieEntries {
			if !entry.IsDir() {
				movies = append(movies, entry.Name())
			}
		}

		// get shows...
		tvEntries, err := os.ReadDir("./media/shows")
		if err != nil {
			fmt.Printf("unable to get show data...")
			http.Error(w, "Unable to read media folder", http.StatusInternalServerError)
			return
		}

		var shows []string
		for _, entry := range tvEntries {
			if entry.IsDir() {
				shows = append(shows, entry.Name())
			}
		}

		// Sort alphabetically
		sort.Strings(movies)

		// Pass to template
		data := PageData{Movies: movies, TVShows: shows}

		tmpl.Execute(w, data)
	})

	// serve movie page route
	// Movie page: show a single video
	http.HandleFunc("/movie/", func(w http.ResponseWriter, r *http.Request) {
		// URL path: /movie/{filename}
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 || parts[2] == "" {
			http.Error(w, "Movie not specified", http.StatusBadRequest)
			return
		}

		title := parts[2]                      // e.g., "movie.mp4"
		ext := filepath.Ext(title)             // ".mp4"
		name := strings.TrimSuffix(title, ext) // "movie"
		ext = strings.TrimPrefix(ext, ".")     // remove the dot -> "mp4"

		data := router.GetMovieInfo(title)
		fmt.Printf("Title: %s\nReleased: %s\nActors: %s\nPlot: %s\n", data.Title, data.Released, data.Actors, data.Plot)

		metaData := router.MetaData{
			Title:    title,
			Released: data.Released,
			Year:     data.Year,
			Rated:    data.Rated,
			Genre:    data.Genre,
			Director: data.Director,
			Actors:   data.Actors,
			Plot:     data.Plot,
			Poster:   data.Poster,
		}

		moviePage.Execute(w, MoviePageData{
			Filename:  name,
			Extension: ext,
			MetaData:  metaData,
		})
	})

	http.HandleFunc("/shows/", func(w http.ResponseWriter, r *http.Request) {
		// URL path:
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 || parts[2] == "" {
			http.Error(w, "Show not specified", http.StatusBadRequest)
			return
		}

		title, err := url.PathUnescape(parts[2])
		if err != nil {
			http.Error(w, "Invalid show name", http.StatusBadRequest)
			return
		}

		data := router.GetSeriesInfo(title)
		//fmt.Printf("Title: %s\nReleased: %s\nActors: %s\nPlot: %s\n", data.Title, data.Released, data.Actors, data.Plot)

		// showName := parts[2]

		metaData := router.MetaData{
			Title:    title,
			Released: data.Released,
			Year:     data.Year,
			Rated:    data.Rated,
			Genre:    data.Genre,
			Director: data.Director,
			Actors:   data.Actors,
			Plot:     data.Plot,
			Poster:   data.Poster,
		}

		showData := ShowPageData{
			MetaData: metaData,
		}

		// iterate through season folders and build season objects with string array of episode file names
		//
		// Base path to the show's folder
		basePath := filepath.Join("./media/shows/", title)

		fmt.Println(basePath)

		// Read all entries in the show's folder
		seasons, err := os.ReadDir(basePath)
		if err != nil {
			http.Error(w, "Error reading show folder", http.StatusInternalServerError)
			return
		}

		for _, season := range seasons {
			if season.IsDir() {
				seasonName := season.Name()
				seasonPath := filepath.Join(basePath, seasonName)

				// Read all episodes in this season
				// will need
				episodeFiles, err := os.ReadDir(seasonPath)
				if err != nil {
					continue // skip this season if there's an error
				}

				var episodes []Episode
				for _, episode := range episodeFiles {
					if !episode.IsDir() {
						// fmt.Println(episode.Name())
						var episodeSingle Episode

						ext := filepath.Ext(episode.Name())
						nameNoExt := strings.TrimSuffix(episode.Name(), ext)
						episodeSingle.Extension = ext
						episodeSingle.FileName = nameNoExt
						episodes = append(episodes, episodeSingle)

						// read subtitle tracks as well if there are any...
						// match via "S01E01" tags...

					}
				}

				// Append season to the showData
				showData.Content = append(showData.Content, SeasonStruct{
					Season:   extractSeasonNumber(seasonName),
					Episodes: episodes,
				})
			}
		}

		fmt.Println(showData)
		showPage.Execute(w, showData)
	})
	// -----------------------------------

	error := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil)
	if error != nil {
		log.Fatal(err)
	}
}
