package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleLogin(w http.ResponseWriter,r *http.Request) {
	conf := &oauth2.Config{
		ClientID:     "GET_THE_ID_FROM_GOOGLE_API_CONSOLE",
		ClientSecret: "GET_CLIENT_SECRET_FROM_GOOGLE_API_CONSOLE",
		RedirectURL:  "http://localhost:8080/callback", // URL it sends the results to, (ie) "Authorization_code"
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state")
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

}

func GoogleCallBack(w http.ResponseWriter,r *http.Request) {
	fmt.Println(r.URL.Query())

	code := r.URL.Query()["code"][0]

	// we got "authorization_code" in the QUERY!
	fmt.Print(code)

	conf := &oauth2.Config{
		ClientID:     "GET_THE_ID_FROM_GOOGLE_API_CONSOLE",
		ClientSecret: "GET_CLIENT_SECRET_FROM_GOOGLE_API_CONSOLE",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // we want email id
		},
		Endpoint: google.Endpoint,
	}

	// Handle the exchange code to initiate a transport.

	// Send POST request containing "authorization_code" to the token Endpoint to get the access token..
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	// tok contains access token

	url1 := "https://www.googleapis.com/oauth2/v3/userinfo?access_token="+ tok.AccessToken;
	client := conf.Client(oauth2.NoContext, tok)
	resp,_ := client.Get(url1)

	str, _ := ioutil.ReadAll(resp.Body)

	// Prints Email Address since we asked for the email address in the scope before :)
	fmt.Print(string(str))
}


func main() {
	mux := http.NewServeMux();

	mux.HandleFunc("/",GoogleLogin)
	mux.HandleFunc("/callback",GoogleCallBack)

	log.Println("started server on :: http://localhost:8080/")
	if oops := http.ListenAndServe(":8080", mux); oops != nil {
		log.Fatal(oops)
	}

}