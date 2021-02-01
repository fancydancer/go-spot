package main

import (
  "fmt"
  "net/http"
  "log"
  //"text/template"
  "html/template"
  "github.com/zmb3/spotify"
  "os"
  "golang.org/x/oauth2/clientcredentials"
  "golang.org/x/oauth2"
  "context"
  //"encoding/json"
)




/* # resources */
/*
* https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.1.html
* https://github.com/zmb3/spotify
*/
/* # resources */

//

/* # http handlers */
func rootHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "<h1>Hello All</h1>")
	t,_ := template.ParseFiles("root.html")
	t.Execute(w, nil)

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles("search.html")
	
	if r.Method == "GET" {
		t.Execute(w, nil)
	}
}

/* # http handlers */

//

/* # spotify Api functions */

func spotSearch(h http.HandlerFunc, token *oauth2.Token) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		t, err := template.ParseFiles("search.html")

	
		if r.Method == "GET" {
			if err != nil {
				fmt.Println(err)
			}
			t.Execute(w, nil)
		} else {
		
			// read input from html
			r.ParseForm()
	
			searchQuery := r.PostFormValue("searchQuery")
	
			// process data
			if searchQuery != "" {
	
				// print the search query
				fmt.Println("Search Query: ", searchQuery)

				client := spotify.Authenticator{}.NewClient(token)
	
				// do da search
				
				results, err := client.Search(searchQuery, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum|spotify.SearchTypeTrack|spotify.SearchTypeArtist)
				//
				var parsedResults = parseSpotifySearchResults(results)


				// new spotifyTrack object
				
				tracks := make([]SpotifyTrack, 0) 
				for _,item := range parsedResults.tracks {
					fmt.Println(item.Name)
					track := SpotifyTrack{item.Name, item.ID.String(), string(item.URI), item.Endpoint}
					tracks = append(tracks, track)
				}
				
				//tracks := parsedResults.tracks
				//albums := parsedResults.albums
				//playlists := parsedResults.playlists
				//artists := parsedResults.artists
				//firstSongName := parsedResults.tracks

				if err != nil {
					fmt.Println("ope")
					log.Println(err)
				}

				t, err := template.ParseFiles("results.html")
				if err != nil {
					fmt.Println(err)
				}
				t.Execute(w, tracks)
	
			} else {
				fmt.Println("Empty search query")
			}
		}
	})
}

//***//

type SpotifySearchResults struct {
	albums []spotify.SimpleAlbum
	playlists []spotify.SimplePlaylist
	tracks []spotify.FullTrack
	artists []spotify.FullArtist
}



type SpotifyTrack struct {
	Name string
	ID string
	URI string
	Endpoint string 
}


/*func (uri *spotify.URI) String() string {
	return string(*)
}*/

/* # custom functions */
func parseSpotifySearchResults(searchResults *spotify.SearchResult) *SpotifySearchResults{

	//spotResults := make([]SpotifySearchResults, 0)
	
	albums := make([]spotify.SimpleAlbum, 0)
	playlists := make([]spotify.SimplePlaylist, 0)
	tracks := make([]spotify.FullTrack, 0)
	artists := make([]spotify.FullArtist, 0)
	
	
	if searchResults.Albums != nil {
		fmt.Println("Albums:")
		for _, item := range searchResults.Albums.Albums {
			fmt.Println("   ", item.Name)	
			albums = append(albums, item)
		}
	}

	if searchResults.Playlists != nil {
		fmt.Println("Playlists:")
		for _, item := range searchResults.Playlists.Playlists {
			fmt.Println("   ", item.Name)
			playlists = append(playlists, item)
		}
	}

	if searchResults.Tracks != nil {
		fmt.Println("Tracks:")
		for _, item := range searchResults.Tracks.Tracks {
			fmt.Println("   ", item.Name)
			tracks = append(tracks, item)
		}
	}

	if searchResults.Artists != nil {
		fmt.Println("Artists:")
		for _, item := range searchResults.Artists.Artists {
			fmt.Println("   ", item.Name)
			artists = append(artists, item)
		}
	}

	spotResults := SpotifySearchResults{albums, playlists, tracks, artists}
	return &spotResults
	


}

/* # pkg main */
func main() {

	// configure the api token request - using "clientcredentials" package
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	// make request for spotify token
	token,err := config.Token(context.Background())
	if err != nil {
		fmt.Println("issue grabbing token")
	}

	// create new spotify api client using the token that was requested
	// spotifyClient = spotify.Authenticator{}.NewClient(token)

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/search", spotSearch(searchHandler, token))

	// results handler
	//http.HandleFunc(/results, searchResults)


	//http.HandleFunc("/results", resultsHandler)
	fmt.Println("listening on 0.0.0.0:8088")
	http.ListenAndServe(":8088", nil)

}