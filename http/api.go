package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/y3sh/go143/twitter"

	"github.com/go-chi/chi"
	"github.com/juju/errors"
)

const (
	SiteRoot     = "/"
	APIRoot      = "/api"
	TweetsURI    = "/api/v1/tweets"
	RandTweetURI = "/api/v1/randTweet"
)

var apiVersion = &APIVersion{"GO143", "v1", []string{
	"/api/v1/tweets",
	"/api/v1/randTweet",
}}

type API struct {
	Router       Router
	TweetService TweetService
}

type Router interface {
	Use(middleware ...func(http.Handler) http.Handler)
	Route(pattern string, fn func(r chi.Router)) chi.Router
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type TweetService interface {
	GetTweets() []*twitter.Tweet
	AddTweet(tweetText string) (*twitter.Tweet, error)
	AddRandTweet() (*twitter.Tweet, error)
}

type APIVersion struct {
	API     string   `json:"api"`
	Version string   `json:"version"`
	URLS    []string `json:"urls"`
}

func NewAPIRouter(httpRouter Router, tweetService TweetService) *API {
	a := &API{
		Router:       httpRouter,
		TweetService: tweetService,
	}

	a.EnableCORS()

	httpRouter.Route(SiteRoot, func(r chi.Router) {
		r.Get("/", a.GetRoot)
	})

	httpRouter.Route(APIRoot, func(r chi.Router) {
		r.Get("/", a.GetTweets)
		r.Post("/", a.PostTweet)
	})

	httpRouter.Route(TweetsURI, func(r chi.Router) {
		r.Get("/", a.GetRandTweet)
	})

	httpRouter.Route(RandTweetURI, func(r chi.Router) {
		r.Get("/", a.GetRandTweet)
	})

	http.Handle(SiteRoot, httpRouter)

	return a
}

func (a *API) GetRoot(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, r, apiVersion)
}

func (a *API) GetTweets(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, r, a.TweetService.GetTweets())
}

func (a *API) PostTweet(w http.ResponseWriter, r *http.Request) {
	var tweet twitter.Tweet

	err := json.NewDecoder(r.Body).Decode(&tweet)
	if err != nil {
		WriteBadRequest(w, r, "Error invalid tweet format.")
	}

	finalTweet, err := a.TweetService.AddTweet(tweet.TweetText)
	if err != nil {
		WriteServerError(w, r, errors.Wrap(err, errors.New("service failed to add tweet")))
	}

	WriteJSON(w, r, finalTweet)
}

func (a *API) GetRandTweet(w http.ResponseWriter, r *http.Request) {
	randTweet, err := a.TweetService.AddRandTweet()
	if err != nil {
		WriteServerError(w, r, errors.Wrap(err, errors.New("service failed to add tweet")))
	}

	WriteJSON(w, r, randTweet)
}

func (a *API) EnableCORS() {
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           100, // Maximum value not ignored by any of major browsers
	})

	a.Router.Use(corsConfig.Handler)
}
