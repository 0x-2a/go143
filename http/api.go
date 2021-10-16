package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/y3sh/go143/instagram"
	"github.com/y3sh/go143/nytimes"
	"github.com/y3sh/go143/repository"
	"github.com/y3sh/go143/twitter"
)

const (
	SiteRoot              = "/"
	TweetsURI             = "/v1/tweets"
	EchoURI               = "/v1/form"
	RandTweetURI          = "/v1/randTweet"
	InstagramUserURI      = "/v1/instagram/users/{cseName}"
	InstagramRandUserURI  = "/v1/instagram/users/random"
	InstagramSessionURI   = "/v1/instagram/sessions/{cseName}"
	NYTimesBestSellersURI = "/v1/nyTimes/bestSellers"
	BookCoverURI          = "/v1/nyTimes/bookCovers/{isbn}"
	FileUploadURI         = "/v1/files"
	ProjectStoreURI       = "/v1/projects/{groupName}/{keyName}"
)

var (
	OK         = &struct{}{}
	apiVersion = &APIVersion{"GO143", "v1.2", []string{
		"https://api.y3sh.com/v1/tweets",
		"https://api.y3sh.com/v1/form",
		"https://api.y3sh.com/v1/randTweet",
		"https://api.y3sh.com/v1/nyTimes/bestSellers",
		"https://api.y3sh.com/v1/nyTimes/bookCovers/{isbn}",
		"https://api.y3sh.com/v1/instagram/users/{cseName}",
		"https://api.y3sh.com/v1/instagram/users/random",
		"https://api.y3sh.com/v1/instagram/session",
		"https://api.y3sh.com/v1/projects/TheATeam/posts",
		"https://api.y3sh.com/v1/files",
	}}
)

type API struct {
	Router               Router
	TweetService         TweetService
	InstagramUserService InstagramUserService
	ProjectStoreService  ProjectStoreService
	S3Repository         S3Repository
	NyTimesClient        NyTimesClient
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
	AddUser(cseName string, user instagram.User) error
	GetUsers(string) []instagram.User
	GetRandProfile() instagram.RandomUser
	IsValidPassword(username string, passwordAttempt string, password string) bool
}

type NyTimesClient interface {
	GetSimpleBestSellers() []nytimes.SimpleBook
	GetBookCoverURL(isbn string) nytimes.BookCoverURL
}

type ProjectStoreService interface {
	GetValue(groupName, keyName string) string
	SetValue(groupName, keyName, value string)
}

type S3Repository interface {
	AddFileToS3(name string, reader *bytes.Reader) (string, error)
}

type APIVersion struct {
	API     string   `json:"api"`
	Version string   `json:"version"`
	URLS    []string `json:"urls"`
}

func NewAPIRouter(httpRouter Router, tweetService TweetService,
	instagramUserService InstagramUserService,
	nyTimesClient NyTimesClient,
	projectStoreService ProjectStoreService,
	s3Repository S3Repository) *API {
	a := &API{
		Router:               httpRouter,
		TweetService:         tweetService,
		InstagramUserService: instagramUserService,
		NyTimesClient:        nyTimesClient,
		ProjectStoreService:  projectStoreService,
		S3Repository:         s3Repository,
	}

	a.EnableCORS()

	httpRouter.Route(SiteRoot, func(r chi.Router) {
		r.Get("/", a.GetRoot)
	})

	httpRouter.Route(TweetsURI, func(r chi.Router) {
		r.Get("/", a.GetTweets)
		r.Post("/", a.PostTweet)
	})

	httpRouter.Route(EchoURI, func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			GetPostEcho(w, r, "get")
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			GetPostEcho(w, r, "post")
		})

		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			GetPostEcho(w, r, "put")
		})

		r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
			GetPostEcho(w, r, "patch")
		})

		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			GetPostEcho(w, r, "delete")
		})
	})

	httpRouter.Route(RandTweetURI, func(r chi.Router) {
		r.Get("/", a.GetRandTweet)
	})

	httpRouter.Route(InstagramUserURI, func(r chi.Router) {
		r.Post("/", a.PostInstagramUser)
		r.Get("/", a.GetInstagramUsers)
	})

	httpRouter.Route(InstagramRandUserURI, func(r chi.Router) {
		r.Get("/", a.GetRandInstagramUser)
	})

	httpRouter.Route(InstagramSessionURI, func(r chi.Router) {
		r.Post("/", a.PostInstagramSession)
	})

	httpRouter.Route(NYTimesBestSellersURI, func(r chi.Router) {
		r.Get("/", a.GetNyTimesBestSellers)
	})

	httpRouter.Route(BookCoverURI, func(r chi.Router) {
		r.Get("/", a.GetNyTimesBookCover)
	})

	httpRouter.Route(ProjectStoreURI, func(r chi.Router) {
		r.Get("/", a.GetProjectKeyValue)
		r.Post("/", a.SetProjectKeyValue)
	})

	httpRouter.Route(FileUploadURI, func(r chi.Router) {
		r.Post("/", a.PostFileUpload)
	})

	http.Handle(SiteRoot, httpRouter)

	return a
}

func GetPostEcho(w http.ResponseWriter, r *http.Request, method string) {
	var q, res string
	if method == "get" {
		q = r.URL.RawQuery
	} else {
		b, _ := ioutil.ReadAll(r.Body)
		q = string(b)
	}

	keyVals := strings.Split(q, "&")
	for _, keyVal := range keyVals {
		leftRight := strings.Split(keyVal, "=")

		if len(leftRight) >= 2 {
			a, _ := url.QueryUnescape(leftRight[0])
			b, _ := url.QueryUnescape(leftRight[1])

			res += fmt.Sprintf("<li>%s: %s</li>", a, b)
		} else if len(leftRight) == 1 {
			a, _ := url.QueryUnescape(leftRight[0])

			res += fmt.Sprintf("<li>%s: </li>", a)
		}
	}

	res = fmt.Sprintf("<ul>%s</ul>", res)

	echoPage := fmt.Sprintf("<!DOCTYPE html><html><body><button class=\"btn\" onclick=\"goBack()\">&larr; Back</button><script>function goBack(){window.history.back();}</script><p>Success! The server got your form.</p><p style=\"margin-left:.5rem;\">Method: %s</p><p style=\"margin-left:.5rem;\">Fields:</p> %s <style>html{width: 100%%; height: 100%%; overflow: hidden;}.btn{padding-left: 1rem; padding-right: 1.3rem; padding-top: .3rem; padding-bottom: .3rem;}body{width: 100%%; height: 100%%; padding-top: 2rem; font-family: 'Helvetica Neue', Arial sans-serif; background: #092756; background: -moz-radial-gradient(0%% 100%%, ellipse cover, rgba(104, 128, 138, .4) 10%%, rgba(138, 114, 76, 0) 40%%), -moz-linear-gradient(top, rgba(57, 173, 219, .25) 0%%, rgba(42, 60, 87, .4) 100%%), -moz-linear-gradient(-45deg, #670d10 0%%, #092756 100%%); background: -webkit-radial-gradient(0%% 100%%, ellipse cover, rgba(104, 128, 138, .4) 10%%, rgba(138, 114, 76, 0) 40%%), -webkit-linear-gradient(top, rgba(57, 173, 219, .25) 0%%, rgba(42, 60, 87, .4) 100%%), -webkit-linear-gradient(-45deg, #670d10 0%%, #092756 100%%); background: -o-radial-gradient(0%% 100%%, ellipse cover, rgba(104, 128, 138, .4) 10%%, rgba(138, 114, 76, 0) 40%%), -o-linear-gradient(top, rgba(57, 173, 219, .25) 0%%, rgba(42, 60, 87, .4) 100%%), -o-linear-gradient(-45deg, #670d10 0%%, #092756 100%%); background: -ms-radial-gradient(0%% 100%%, ellipse cover, rgba(104, 128, 138, .4) 10%%, rgba(138, 114, 76, 0) 40%%), -ms-linear-gradient(top, rgba(57, 173, 219, .25) 0%%, rgba(42, 60, 87, .4) 100%%), -ms-linear-gradient(-45deg, #670d10 0%%, #092756 100%%); background: -webkit-radial-gradient(0%% 100%%, ellipse cover, rgba(104, 128, 138, .4) 10%%, rgba(138, 114, 76, 0) 40%%), linear-gradient(to bottom, rgba(57, 173, 219, .25) 0%%, rgba(42, 60, 87, .4) 100%%), linear-gradient(135deg, #670d10 0%%, #092756 100%%); filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#3E1D6D', endColorstr='#092756', GradientType=1);}body{margin-left: 2rem;}li, h4, p{color: #fff; list-style: none; font-size: 1.5rem;}</style></body></html>", method, res)

	WriteResponse(w, r, []byte(echoPage))

}

func (a *API) PostInstagramUser(w http.ResponseWriter, r *http.Request) {
	cseName := chi.URLParam(r, "cseName")
	if cseName == "" {
		WriteBadRequest(w, r, "Missing CSE Name")
		return
	}

	var user instagram.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteBadRequest(w, r, "Error invalid user format.")
		return
	}

	err = a.InstagramUserService.AddUser(cseName, user)
	if err != nil {
		WriteBadRequest(w, r, fmt.Sprintf("Error: %s.", err.Error()))
		return
	}

	WriteJSON(w, r, OK)
}

func (a *API) GetInstagramUsers(w http.ResponseWriter, r *http.Request) {
	cseName := chi.URLParam(r, "cseName")
	if cseName == "" {
		WriteBadRequest(w, r, "Missing CSE Name")
		return
	}

	users := a.InstagramUserService.GetUsers(cseName)
	WriteJSON(w, r, users)
}

func (a *API) PostInstagramSession(w http.ResponseWriter, r *http.Request) {
	cseName := chi.URLParam(r, "cseName")
	if cseName == "" {
		WriteBadRequest(w, r, "Missing CSE Name")
		return
	}

	var user instagram.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteBadRequest(w, r, "Error invalid user format.")
		return
	}

	validPassword := a.InstagramUserService.IsValidPassword(cseName, user.Username, user.Password)
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

func (a *API) GetNyTimesBestSellers(w http.ResponseWriter, r *http.Request) {
	bestSellers := a.NyTimesClient.GetSimpleBestSellers()

	WriteJSON(w, r, bestSellers)
}

func (a *API) GetNyTimesBookCover(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	coverURL := a.NyTimesClient.GetBookCoverURL(isbn)

	WriteJSON(w, r, coverURL)
}

func (a *API) GetProjectKeyValue(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "groupName")
	keyName := chi.URLParam(r, "keyName")

	val := a.ProjectStoreService.GetValue(groupName, keyName)

	w.Header().Set("content-type", "application/json")

	if val == "" {
		WriteResponse(w, r, []byte("null"))
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

func (a *API) PostFileUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(256 * 1000) // bytes

	var buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Errorf("Could not read file \n%+v\n", err)

		WriteBadRequest(w, r, "Could not read file.")
		return
	}

	nameParts := strings.Split(header.Filename, ".")
	name := uuid.New().String()

	if len(nameParts) > 1 {
		name = fmt.Sprintf("%s.%s", name, nameParts[len(nameParts)-1])
	}

	// Copy the file data to my buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		WriteError(w, r, "Err buffering file", 500)
		return
	}

	fileURL, err := a.S3Repository.AddFileToS3(name, bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Errorf("%+v", err)
		WriteError(w, r, "Err uploading to S3", 500)
		return
	}

	WriteJSON(w, r, repository.S3Response{
		FileURL: fileURL,
	})
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

func (a *API) GetRandInstagramUser(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, r, a.InstagramUserService.GetRandProfile())
}
