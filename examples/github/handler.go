package main

import (
	"fmt"
	"net/http"
)

const (
	HelloTemp = `
Hello %s, you have awesome %d public repos
	`
)

type Handler struct {
	Config      Config
	Client      *Client
	ProfileRepo ProfileRepo
}

func NewHandler(cfg Config, client *Client, profileRepo ProfileRepo) Handler {
	return Handler{
		Config:      cfg,
		Client:      client,
		ProfileRepo: profileRepo,
	}
}

func (h Handler) Home(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Welcome"))
}

func (h Handler) Health(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Pong!"))
}

func (h Handler) Auth(w http.ResponseWriter, req *http.Request) {
	url, err := h.Client.GetAuthorizationURL()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, req, url, 301)
}

func (h Handler) Callback(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	codes, exists := queryValues["code"]
	if !exists {
		w.Write([]byte("Code couldn't be found"))
		return
	}
	if len(codes) == 0 {
		w.Write([]byte("Code couldn't be found"))
		return
	}
	res, err := h.Client.GetAccessToken(codes[0])
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	profile, err := h.Client.GetAuthenticated(res.String())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if err := h.ProfileRepo.Store(profile); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	profile, err = h.ProfileRepo.ResolveByID(profile.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf(HelloTemp, profile.Name, profile.PublicRepos)))
}
