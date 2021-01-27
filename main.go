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
  "context"
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
	t,_ := template.ParseFiles("root.gtpl")
	t.Execute(w, nil)

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles("search.gtpl")
	
	if r.Method == "GET" {
		t.Execute(w, nil)
	}
}

/* # http handlers */

//

/* # spotify Api functions */

func spotSearch(h http.HandlerFunc, client spotify.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		t,_ := template.ParseFiles("search.gtpl")
	
		if r.Method == "GET" {
			t.Execute(w, nil)
		} else {
		
			// read input from html
			r.ParseForm()
	
			searchQuery := r.PostFormValue("searchQuery")
	
			// process data
			if searchQuery != "" {
	
				// print the search query
				fmt.Println("Search Query: ", searchQuery)
	
				// do da search
				
				results, err := client.Search(searchQuery, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum|spotify.SearchTypeTrack|spotify.SearchTypeArtist)
				//
	
	
				if err != nil {
					log.Println(err)
				}
	
				// check for albums
				if results.Albums != nil {
					fmt.Println("Albums:")
					for _, item := range results.Albums.Albums {
						fmt.Println("   ", item.Name)
					}
				}
	
				//check for playlists
				if results.Playlists != nil {
					fmt.Println("Playlists:")
					for _, item := range results.Playlists.Playlists {
						fmt.Println("   ", item.Name)
					}
				}
	
				// check for songs
				if results.Tracks != nil {
					fmt.Println("Tracks:")
					for _, item := range results.Tracks.Tracks {
						fmt.Println("   ", item.Name)
					}
				}
	
				//check for artists
				if results.Artists != nil {
					fmt.Println("Artists:")
					for _, item := range results.Artists.Artists {
						fmt.Println("   ", item.Name)
					}
				}
				
				searchResults := "some stuff"
				fmt.Println(searchResults)
	
				t,_ := template.ParseFiles("results.gtpl")
				t.Execute(w, nil)
	
			} else {
				fmt.Println("Empty search query")
			}
			t.Execute(w, nil)
		}
	})
}

/* # spotify API functions */

//

// redirectURL is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURL = "http://localhost:8088/callback/"

/* pkg main */
func main() {

	var spotifyClient spotify.Client
	
	fmt.Println(spotify.TokenURL)
	
	// configure the api token request - using "clientcredentials" package
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	// make request for spotify token
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}
	else {
		fmt.Println("Successfully obtained API token for Spotify")
	}


	spotifyClient := spotify.Authenticator{}.NewClient(token)
	
	// web request handlers
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		// set spotify client id and secret variables
		spotify_id := os.Getenv("SPOTIFY_ID")
		
		spotify_secret := os.Getenv("SPOTIFY_SECRET")
		
		// create authenticator for new spotify client
		auth := spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate)
		auth.SetAuthInfo(spotify_id, spotify_secret)
		fmt.Println(spotify_id, spotify_secret)
		state := "abc123"
		url := auth.AuthURL(state)
		fmt.Println(url)
		token,err := auth.Token(state, r)
		if err != nil {
			fmt.Println("There was an error with grabbing the authentiction token")
		}
		fmt.Println(token)
		// grab authentication credentials for spotify client from OS environment variables
		spotifyClient = auth.NewClient(token)
	})
	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/search", spotSearch(searchHandler, spotifyClient))

	// results handler
	//http.HandleFunc(/results, searchResults)


	//http.HandleFunc("/results", resultsHandler)
	fmt.Println("listening on 0.0.0.0:8088")
	http.ListenAndServe(":8088", nil)

}