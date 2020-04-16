package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2/jwt"
)

// AuthorizeUser implements the PKCE OAuth2 flow.
func AuthorizeUser(clientID string, authDomain string, redirectURL string) {
	// initialize the code verifier
	var CodeVerifier, _ = cv.CreateCodeVerifier()

	// Create code_challenge with S256 method
	codeChallenge := CodeVerifier.CodeChallengeS256()

	// construct the authorization URL (with Auth0 as the authorization provider)
	authorizationURL := fmt.Sprintf(
		"https://%s/authorize?audience=https://api.snapmaster.io"+
			"&scope=openid+profile+email"+
			"&response_type=code&client_id=%s"+
			"&code_challenge=%s"+
			"&code_challenge_method=S256&redirect_uri=%s",
		authDomain, clientID, codeChallenge, redirectURL)

	// start a web server to listen on a callback URL
	server := &http.Server{Addr: redirectURL}

	// define a handler that will get the authorization code, call the token endpoint, and close the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get the authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Println("snap: Url Param 'code' is missing")
			io.WriteString(w, "Error: could not find 'code' URL parameter\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// trade the authorization code and the code verifier for an access token
		codeVerifier := CodeVerifier.String()
		responseData, err := getAccessToken(clientID, codeVerifier, code, redirectURL)
		if err != nil {
			fmt.Println("snap: could not get access token")
			io.WriteString(w, "Error: could not retrieve access token\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// retrieve the idtoken and access token out of the map, and return to caller
		idToken := responseData["id_token"].(string)
		accessToken := responseData["access_token"].(string)

		// parse the id_token JWT into its claims
		claims := parseJWT(idToken)
		name := claims["name"].(string)
		email := claims["email"].(string)

		// store some user identity claims
		viper.Set("Name", name)
		viper.Set("Email", email)

		// set the access token
		viper.Set("AccessToken", accessToken)

		// store the config
		err = viper.WriteConfig()
		//_, err = config.WriteConfigFile("auth.json", token)
		if err != nil {
			fmt.Println("snap: could not write config file")
			io.WriteString(w, "Error: could not store access token\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// return an indication of success to the caller
		io.WriteString(w, `
		<html>
			<head>
				<link href="https://fonts.googleapis.com/css?family=Lato:100,300,400&display=swap" rel="stylesheet">  
				<title>SnapMaster</title>
			</head>
			<body style="background: #000; color: #fff; font-family: 'Lato', -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen",
			"Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
			sans-serif; font-weight: 300;">
				<center style="margin: 100">
					<img src="https://www.snapmaster.io/SnapMaster-logo-220.png" alt="snapmaster" width="100" height="100" />
					<h1>Hi, `)
		io.WriteString(w, name)
		io.WriteString(w, `! You've logged in successfully.</h1>
					<h2>You can close this window and return to the snap CLI.</h2>
				</center>
			</body>
		</html>`)

		fmt.Println("Successfully logged into snapmaster API.")

		// close the HTTP server
		cleanup(server)
	})

	// parse the redirect URL for the port number
	u, err := url.Parse(redirectURL)
	if err != nil {
		fmt.Printf("snap: bad redirect URL: %s\n", err)
		os.Exit(1)
	}

	// set up a listener on the redirect port
	port := fmt.Sprintf(":%s", u.Port())
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("snap: can't listen to port %s: %s\n", port, err)
		os.Exit(1)
	}

	// open a browser window to the authorizationURL
	err = open.Start(authorizationURL)
	if err != nil {
		fmt.Printf("snap: can't open browser to URL %s: %s\n", authorizationURL, err)
		os.Exit(1)
	}

	// start the blocking web server loop
	// this will exit when the handler gets fired and calls server.Close()
	server.Serve(l)
}

// getAccessToken trades the authorization code retrieved from the first OAuth2 leg for an access token
func getAccessToken(clientID string, codeVerifier string, authorizationCode string, callbackURL string) (map[string]interface{}, error) {
	// set the url and form-encoded data for the POST to the access token endpoint
	url := "https://snapmaster-dev.auth0.com/oauth/token"
	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s"+
			"&code_verifier=%s"+
			"&code=%s"+
			"&redirect_uri=%s",
		clientID, codeVerifier, authorizationCode, callbackURL)
	payload := strings.NewReader(data)

	// create the request and execute it
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("snap: HTTP error: %s", err)
		return nil, err
	}

	// process the response
	defer res.Body.Close()
	var responseData map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)

	// unmarshal the json into a string map
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("snap: JSON error: %s", err)
		return nil, err
	}

	return responseData, nil
}

// cleanup closes the HTTP server
func cleanup(server *http.Server) {
	// we run this as a goroutine so that this function falls through and
	// the socket to the browser gets flushed/closed before the server goes away
	go server.Close()
}

func parseJWT(tokenString string) map[string]interface{} {
	var claims map[string]interface{} // generic map to store parsed token

	// decode JWT token without verifying the signature
	token, err := jwt.ParseSigned(tokenString)
	if err != nil {
		fmt.Printf("snap: could not parse JWT\nerror: %s\n", err)
		os.Exit(1)
	}

	err = token.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		fmt.Printf("snap: could not parse JWT\nerror: %s\n", err)
		os.Exit(1)
	}

	/*
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("<YOUR VERIFICATION KEY>"), nil
		})

		if err != nil {
			fmt.Printf("snap: could not parse JWT\nerror: %s\n", err)
			os.Exit(1)
		}
	*/
	return claims
	/*
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims
		}

		fmt.Printf("snap: could not parse JWT\nerror: %s\n", err)
		os.Exit(1)
		return nil*/
}
