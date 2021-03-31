package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var httpClient http.Client

const (
	//github
	//ClientID     string = "02d37221f52028d30d08"
	//ClientSecret string = "7987abeefde0c400edafb9892b82c7566a9d2e96"

	//OpenAM
	clientID     string = "myTestAuth"
	clientSecret string = "83WtpaKeQb9aY6B"
)

func init() {
	// We will be using `httpClient` to make external HTTP requests later in our code
	httpClient = http.Client{}
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/oauth/redirect", handleRedirect)
	http.ListenAndServe(":8081", nil)
}

// Create a new redirect route route
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// First, we need to get the value of the `code` query param
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	fmt.Printf("r: %+v\n", r)
	code := r.FormValue("code")

	// Next, lets for the HTTP request to call the github oauth enpoint
	// to get our access token
	//reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	reqURL := fmt.Sprintf("http://localhost:8080/opnam/oauth2/realms/root/access_token?grant_type=client_credentials&client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	fmt.Println(reqURL)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	// We set this header since we want the response
	// as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	var buf []byte
	res.Body.Read(buf)
	fmt.Printf("response: %v\n", string(buf))

	// Parse the request body into the `OAuthAccessResponse` struct
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// Finally, send a response to redirect the user to the "welcome" page
	// with the access token
	w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
	w.WriteHeader(http.StatusFound)
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
