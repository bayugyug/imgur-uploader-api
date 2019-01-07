package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var (
	pContext  context.Context
	pAuthConf *oauth2.Config
	pAuthURL  string
)

//initAuthBearerToken try to init the token grabber
func initAuthBearerToken() string {

	pContext = context.Background()
	pAuthConf = &oauth2.Config{
		ClientID:     pParamConfig.ID,
		ClientSecret: pParamConfig.Secret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			TokenURL: pImgurTokenURL,
			AuthURL:  pImgurAuthURL,
		},
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	pAuthURL = pAuthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	return pAuthURL
}

//getAuthBearerToken et the auth token based on the code passed as params
func getAuthBearerToken() string {

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.

	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 30 * time.Second}
	pContext = context.WithValue(pContext, oauth2.HTTPClient, httpClient)

	// shake it
	tok, err := pAuthConf.Exchange(pContext, pParamConfig.Code)
	if err != nil {
		log.Println(err)
		if pParamConfig.Bearer != nil && pParamConfig.Bearer.RefreshToken != "" {
			_ = checkOldBearerToken(err.Error())
		}
		return err.Error()
	}

	// grab or refresh
	client := pAuthConf.Client(pContext, tok)
	_ = client

	//save the bearer
	pParamConfig.Bearer = tok

	dumper(pParamConfig)

	//good
	return ""
}

//showMacroMsg
func showMacroMsg(url string) {

	//fmt msg
	msg := `

	Visit the URL for the auth dialog: (use Chrome and approve it)

	` + url + `

	@Replace the code in the Macro below: 

	    {CODE_FROM_THE_REDIRECT_URL}

	curl -v -X GET  'http://127.0.0.1:` + pHttpPort + `/v1/api/credentials/{CODE_FROM_THE_REDIRECT_URL}'

	`

	fmt.Println(msg)
}

//fmtBearerHeader set token in header
func fmtBearerHeader() map[string]string {
	//Authorization: Bearer {refresh_token}
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", pParamConfig.Bearer.AccessToken),
	}
}

//fmtBearerHeaderRefreshToken set token in header
func fmtBearerHeaderRefreshToken() (map[string]string, string) {
	frm := url.Values{}
	frm.Set("client_id", pParamConfig.ID)
	frm.Set("client_secret", pParamConfig.Secret)
	frm.Set("refresh_token", pParamConfig.Bearer.RefreshToken)
	frm.Set("grant_type", "refresh_token")
	return map[string]string{
		"Authorization": fmt.Sprintf("Client-ID: %s", pParamConfig.ID),
	}, frm.Encode()
}

//checkOldBearerToken check if previously
func checkOldBearerToken(errStr string) *oauth2.Token {

	if strings.Contains(errStr, "400 Bad Request") &&
		strings.Contains(errStr, "The authorization code has expired") {
		//try to get the refresh
		hdrs, frm := fmtBearerHeaderRefreshToken()
		postresp, postcode, perr := httpPost(pImgurTokenURL, frm, hdrs)
		if postcode != http.StatusOK || perr != nil {
			return nil //Oops
		}
		dumper(postresp)
	}
	return nil //Oops
}
