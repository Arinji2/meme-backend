package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Arinji2/meme-backend/api"
	custom_log "github.com/Arinji2/meme-backend/logger"
	"github.com/Arinji2/meme-backend/types"
)

func RegisterWithGoogleOauth(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	if authCode == "" {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		http.Error(w, "Missing Environment Variables", http.StatusInternalServerError)
		return
	}

	client := api.NewApiClient("https://oauth2.googleapis.com")

	result, status, err := client.SendRequestWithBody("POST", "/token", map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  redirectURI,
		"grant_type":    "authorization_code",
		"code":          authCode,
	}, map[string]string{}, "application/x-www-form-urlencoded")

	if err != nil || status != 200 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	refreshToken, ok := result["refresh_token"].(string)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	accessToken, err := GetGoogleAccessToken(refreshToken)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		custom_log.Logger.Errorf("Error Getting Google Access Token: %v", err)
		return
	}
	userInfo, err := getGoogleUserInfo(accessToken)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		custom_log.Logger.Errorf("Error Getting Google User Info: %v", err)
		return
	}
	fmt.Println(userInfo.Email)

}

func GetGoogleAccessToken(refreshToken string) (string, error) {
	client := api.NewApiClient("https://oauth2.googleapis.com")

	result, status, err := client.SendRequestWithBody("POST", "/token", map[string]string{
		"client_id":     os.Getenv("GOOGLE_CLIENT_ID"),
		"client_secret": os.Getenv("GOOGLE_CLIENT_SECRET"),
		"redirect_uri":  os.Getenv("GOOGLE_REDIRECT_URI"),
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}, map[string]string{}, "application/x-www-form-urlencoded")

	if err != nil {

		return "", errors.New("at fetching data" + err.Error())
	}

	if status != 200 {

		return "", errors.New("at status code: " + string(status))
	}

	return result["access_token"].(string), nil

}

func getGoogleUserInfo(accessToken string) (types.GoogleUserInfo, error) {
	client := api.NewApiClient("https://www.googleapis.com/oauth2/v3")

	result, status, err := client.SendRequestWithBody("GET", "/userinfo", map[string]string{}, map[string]string{
		"Authorization": "Bearer " + accessToken,
	}, "application/x-www-form-urlencoded")

	if err != nil {

		return types.GoogleUserInfo{}, errors.New("at fetching data" + err.Error())
	}
	if status != 200 {
		fmt.Println(status)
		return types.GoogleUserInfo{}, errors.New("at fetching data, status code: " + string(status))
	}

	jsonResult, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		return types.GoogleUserInfo{}, errors.New("at marshalling, json error: " + jsonErr.Error())
	}

	var googleUserInfo types.GoogleUserInfo
	jsonErr = json.Unmarshal(jsonResult, &googleUserInfo)
	if jsonErr != nil {

		return types.GoogleUserInfo{}, errors.New("at unmarshalling, json error: " + jsonErr.Error())
	}

	return googleUserInfo, nil

}

/*
https://accounts.google.com/o/oauth2/v2/auth
?client_id=413510440293-ukeo9nf35h3hdg9j72f74sfpr8i910ci.apps.googleusercontent.com
&redirect_uri=http://localhost:8080/oauth2-redirect/google
&response_type=code
&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email%20https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile%20openid
&access_type=offline
*/
