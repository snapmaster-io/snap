package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
)

var CodeVerifier, _ = cv.CreateCodeVerifier()

func GetUserAuthorization(clientId string, redirectUri string) string {
	// Create code_verifier
	//v, _ := go_oauth_pkce_code_verifier.CreateCodeVerifier()
	//codeVerifier := v.String()

	// Create code_challenge with S256 method
	//codeChallenge := v.CodeChallengeS256()
	codeChallenge := CodeVerifier.CodeChallengeS256()

	url := "https://snapmaster-dev.auth0.com/authorize?audience=https://api.snapmaster.io&" +
		//	scope=SCOPE&
		"response_type=code&client_id=" + clientId +
		"&code_challenge=" + codeChallenge +
		"&code_challenge_method=S256&redirect_uri=" + redirectUri

	return url
}

func GetAccessToken(clientId string, authorizationCode string, callbackUri string) (string, error) {

	codeVerifier := CodeVerifier.String()

	url := "https://snapmaster-dev.auth0.com/oauth/token"

	data := "grant_type=authorization_code&client_id=" + clientId + "&code_verifier=" + codeVerifier + "&code=" + authorizationCode + "&redirect_uri=" + callbackUri
	//payload := strings.NewReader("grant_type=authorization_code&client_id=%24%7Baccount.clientId%7D&code_verifier=YOUR_GENERATED_CODE_VERIFIER&code=YOUR_AUTHORIZATION_CODE&redirect_uri=%24%7Baccount.callback%7D")
	payload := strings.NewReader(data)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	return string(body), nil
}
