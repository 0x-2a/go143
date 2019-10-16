package http

import (
	"encoding/json"
	"net/http"

	"github.com/y3sh/go143/twitter"

	"github.com/go-chi/chi"
	"github.com/juju/errors"
)

type APIRouter struct {
	TweetService TweetService
}

type TweetService interface {
	GetTweets() []*twitter.Tweet
	AddTweet(tweet *twitter.Tweet) error
}

type APIVersion struct {
	API     string   `json:"api"`
	Version string   `json:"version"`
	URLS    []string `json:"urls"`
}

var apiVersion = &APIVersion{"GO143", "v1", []string{
	"/api/v1/tweets",
}}

func NewAPIRouter(tweetService TweetService) *APIRouter {
	return &APIRouter{
		TweetService: tweetService,
	}
}

func (a *APIRouter) HandleRoot(r chi.Router) {
	r.Get("/", a.GetRoot)
}

func (a *APIRouter) HandleTweets(r chi.Router) {
	r.Get("/", a.GetTweets)
	r.Post("/", a.PostTweet)
}

func (a *APIRouter) GetRoot(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, r, apiVersion)
}

func (a *APIRouter) GetTweets(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, r, a.TweetService.GetTweets())
}

func (a *APIRouter) PostTweet(w http.ResponseWriter, r *http.Request) {
	var tweet twitter.Tweet

	err := json.NewDecoder(r.Body).Decode(&tweet)
	if err != nil {
		WriteBadRequest(w, r, "Error invalid tweet format.")
	}

	err = a.TweetService.AddTweet(&tweet)
	if err != nil {
		WriteServerError(w, r, errors.Wrap(err, errors.New("service failed to add tweet")))
	}

	WriteJSON(w, r, "")
}
