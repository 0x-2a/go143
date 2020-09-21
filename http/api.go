package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/juju/errors"
	"github.com/y3sh/go143/instagram"
	"github.com/y3sh/go143/twitter"
)

const (
	SiteRoot            = "/"
	APIRoot             = "/api"
	TweetsURI           = "/api/v1/tweets"
	RandTweetURI        = "/api/v1/randTweet"
	InstagramUserURI    = "/api/v1/instagram/users"
	InstagramSessionURI = "/api/v1/instagram/sessions"
	ProjectStoreURI     = "/api/v1/projects/{groupName}/{keyName}"
)

var (
	OK         = &struct{}{}
	apiVersion = &APIVersion{"GO143", "v1", []string{
		"https://cos143xl.cse.taylor.edu:8080/api/v1/tweets",
		"https://cos143xl.cse.taylor.edu:8080/api/v1/randTweet",
		"https://cos143xl.cse.taylor.edu:8080/api/v1/instagram/user",
		"https://cos143xl.cse.taylor.edu:8080/api/v1/instagram/session",
		"https://cos143xl.cse.taylor.edu:8080/api/v1/projects/TheATeam/posts",
	}}
)

type API struct {
	Router               Router
	TweetService         TweetService
	InstagramUserService InstagramUserService
	ProjectStoreService  ProjectStoreService
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

type InstagramUserService interface {
	AddUser(user instagram.User) error
	IsValidPassword(username instagram.Username, passwordAttempt string) bool
}

type ProjectStoreService interface {
	GetValue(groupName, keyName string) string
	SetValue(groupName, keyName, value string)
}

type APIVersion struct {
	API     string   `json:"api"`
	Version string   `json:"version"`
	URLS    []string `json:"urls"`
}

func NewAPIRouter(httpRouter Router, tweetService TweetService, instagramUserService InstagramUserService,
	projectStoreService ProjectStoreService) *API {
	a := &API{
		Router:               httpRouter,
		TweetService:         tweetService,
		InstagramUserService: instagramUserService,
		ProjectStoreService:  projectStoreService,
	}

	a.EnableCORS()

	httpRouter.Route(SiteRoot, func(r chi.Router) {
		r.Get("/", a.GetRoot)
	})

	httpRouter.Route(APIRoot, func(r chi.Router) {
		r.Get("/", a.GetRoot)
	})

	httpRouter.Route(TweetsURI, func(r chi.Router) {
		r.Get("/", a.GetTweets)
		r.Post("/", a.PostTweet)
	})

	httpRouter.Route(RandTweetURI, func(r chi.Router) {
		r.Get("/", a.GetRandTweet)
	})

	httpRouter.Route(InstagramUserURI, func(r chi.Router) {
		r.Post("/", a.PostInstagramUser)
	})

	httpRouter.Route(InstagramSessionURI, func(r chi.Router) {
		r.Post("/", a.PostInstagramSession)
	})

	httpRouter.Route(ProjectStoreURI, func(r chi.Router) {
		r.Get("/", a.GetProjectKeyValue)
		r.Post("/", a.SetProjectKeyValue)
	})

	http.Handle(SiteRoot, httpRouter)

	return a
}

func (a *API) PostInstagramUser(w http.ResponseWriter, r *http.Request) {
	var user instagram.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteBadRequest(w, r, "Error invalid user format.")
		return
	}

	err = a.InstagramUserService.AddUser(user)
	if err != nil {
		WriteBadRequest(w, r, fmt.Sprintf("Error: %s.", err.Error()))
		return
	}

	WriteJSON(w, r, OK)
}

func (a *API) PostInstagramSession(w http.ResponseWriter, r *http.Request) {
	var user instagram.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteBadRequest(w, r, "Error invalid user format.")
		return
	}

	validPassword := a.InstagramUserService.IsValidPassword(user.Username, user.Password)
	if !validPassword {
		WriteBadRequest(w, r, "Invalid username and/or password.")
		return
	}

	WriteJSON(w, r, OK)
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
		return
	}

	tweetLen := len(tweet.TweetText)
	if tweetLen < 1 || tweetLen > 280 {
		WriteBadRequest(w, r, "Tweet length must be 1-280 characters.")
		return
	}

	finalTweet, err := a.TweetService.AddTweet(tweet.TweetText)
	if err != nil {
		WriteServerError(w, r, errors.Wrap(err, errors.New("service failed to add tweet")))
		return
	}

	WriteJSON(w, r, finalTweet)
}

func (a *API) GetRandTweet(w http.ResponseWriter, r *http.Request) {
	randTweet, err := a.TweetService.AddRandTweet()
	if err != nil {
		WriteServerError(w, r, errors.Wrap(err, errors.New("service failed to add tweet")))
		return
	}

	WriteJSON(w, r, randTweet)
}

func (a *API) GetProjectKeyValue(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "groupName")
	keyName := chi.URLParam(r, "keyName")

	val := a.ProjectStoreService.GetValue(groupName, keyName)

	w.Header().Set("content-type", "application/json")

	if val == "" {
		WriteError(w, r, "No data found, try a post first.", http.StatusBadRequest)
		return
	}

	WriteResponse(w, r, []byte(val))
}

func (a *API) SetProjectKeyValue(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteBadRequest(w, r, "Error invalid request body.")
		return
	}

	var ignored interface{}

	err = json.Unmarshal(bodyBytes, &ignored)
	if err != nil {
		WriteBadRequest(w, r, "Error invalid JSON.")
		return
	}

	groupName := chi.URLParam(r, "groupName")
	keyName := chi.URLParam(r, "keyName")

	a.ProjectStoreService.SetValue(groupName, keyName, string(bodyBytes))

	WriteJSON(w, r, OK)
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
