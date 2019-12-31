package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

type oidcData struct {
	AuthCode  string
	AuthToken *oauth2.Token
	UserData  jwt.Claims
}

// VerificationData : This struct stores the API response to verify call
type VerificationData struct {
	Success bool
	Msg     string
}

func main() {
	var (
		serverPort        = "3000"
		serverURL         = "http://localhost:" + serverPort
		keyCloakServerURL = "<<Enter the OIDC client link>>"
		clientID          = "golang-client"
		redirectURL       = serverURL + "/oidc/callback"
	)

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, keyCloakServerURL)
	if err != nil {
		log.Fatalf(err.Error())
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	oidcVerifier := provider.Verifier(oidcConfig)

	oauth2Config := &oauth2.Config{
		ClientID:    clientID,
		Endpoint:    provider.Endpoint(),
		RedirectURL: redirectURL,
		Scopes:      []string{oidc.ScopeOpenID},
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Initialize with default dummy values
	reqData := &oidcData{}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.Must(template.ParseFiles("./assets/index.html")).Clone()
		t.Execute(w, reqData)
		// http.Redirect(w, r, "/oidc/login", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("/oidc/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, oauth2Config.AuthCodeURL("state"), http.StatusTemporaryRedirect)
	})

	mux.HandleFunc("/oidc/callback", func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		authCode := queryValues.Get("code")
		authToken, err := oauth2Config.Exchange(ctx, authCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reqData.AuthCode = authCode
		reqData.AuthToken = authToken
		reqData.UserData = parseJWT(authToken.AccessToken)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("/oidc/verify", func(w http.ResponseWriter, r *http.Request) {
		data := &VerificationData{
			Success: false,
			Msg:     "",
		}

		rawIDToken, ok := reqData.AuthToken.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}
		IDToken, err := oidcVerifier.Verify(ctx, rawIDToken)
		if err != nil {
			data.Msg = err.Error()
		} else {
			data.Success = true
			IDTokenJSON, _ := json.Marshal(IDToken)
			data.Msg = string(IDTokenJSON)
		}
		jsonData, encodeErr := json.Marshal(data)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonData)
	})

	mux.HandleFunc("/oidc/refresh", func(w http.ResponseWriter, r *http.Request) {
		tokenSource := oauth2Config.TokenSource(ctx, reqData.AuthToken)
		newToken, refreshErr := tokenSource.Token()
		if refreshErr != nil {
			http.Error(w, refreshErr.Error(), http.StatusInternalServerError)
			return
		}
		reqData.AuthToken = newToken
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	fmt.Println("Attempting to start the server at 0.0.0.0:8000")
	log.Fatalf(http.ListenAndServe(":"+serverPort, mux).Error())
}

func parseJWT(ts string) jwt.Claims {
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(ts, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	return token.Claims
}
