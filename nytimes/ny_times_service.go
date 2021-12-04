package nytimes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type BestSellerRes struct {
	Status       string   `json:"status"`
	Copyright    string   `json:"copyright"`
	NumResults   int64    `json:"num_results"`
	LastModified string   `json:"last_modified"`
	Results      []Result `json:"results"`
}

type Result struct {
	ListName         string       `json:"list_name"`
	DisplayName      string       `json:"display_name"`
	BestsellersDate  string       `json:"bestsellers_date"`
	PublishedDate    string       `json:"published_date"`
	Rank             int64        `json:"rank"`
	RankLastWeek     int64        `json:"rank_last_week"`
	WeeksOnList      int64        `json:"weeks_on_list"`
	Asterisk         int64        `json:"asterisk"`
	Dagger           int64        `json:"dagger"`
	AmazonProductURL string       `json:"amazon_product_url"`
	Isbns            []Isbn       `json:"isbns"`
	BookDetails      []BookDetail `json:"book_details"`
	Reviews          []Review     `json:"reviews"`
}

type BookDetail struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Contributor     string `json:"contributor"`
	Author          string `json:"author"`
	ContributorNote string `json:"contributor_note"`
	Price           string `json:"price"`
	AgeGroup        string `json:"age_group"`
	Publisher       string `json:"publisher"`
	PrimaryIsbn13   string `json:"primary_isbn13"`
	PrimaryIsbn10   string `json:"primary_isbn10"`
}

type Isbn struct {
	Isbn10 string `json:"isbn10"`
	Isbn13 string `json:"isbn13"`
}

type Review struct {
	BookReviewLink     string `json:"book_review_link"`
	FirstChapterLink   string `json:"first_chapter_link"`
	SundayReviewLink   string `json:"sunday_review_link"`
	ArticleChapterLink string `json:"article_chapter_link"`
}

type SimpleBook struct {
	LastWeekRank     int64  `json:"lastWeekRank"`
	WeeksOnList      int64  `json:"weeksOnList"`
	Rank             int64  `json:"rank"`
	AmazonProductURL string `json:"amazonProductUrl"`
	Publisher        string `json:"publisher"`
	Title            string `json:"title"`
	Author           string `json:"author"`
	Description      string `json:"description"`
	Isbn             string `json:"isbn"`
}

type BookCoverURL struct {
	URL string `json:"url"`
}

type GoogleBookRes struct {
	Items []struct {
		VolumeInfo struct {
			ImageLinks struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (resp *http.Response, err error)
}

type RestClient struct {
	bestSellerAPIKey string
	googleBookAPIKey string
	lastBestSellers  BestSellerRes
	isbnCoverURLMap  sync.Map // not thread safe
	httpClient       HTTPClient
}

func NewRestClient(bestSellerAPIKey, googleBookAPIKey string, httpClient HTTPClient) *RestClient {
	return &RestClient{
		bestSellerAPIKey: bestSellerAPIKey,
		googleBookAPIKey: googleBookAPIKey,
		isbnCoverURLMap:  sync.Map{},
		httpClient:       httpClient,
	}
}

func (r *RestClient) GetSimpleBestSellers() []SimpleBook {
	var books []SimpleBook

	bestSellers := r.getBestSellers()
	for i := range bestSellers.Results {
		bookObj := bestSellers.Results[i]
		if len(bookObj.BookDetails) > 0 {
			bookInfo := bookObj.BookDetails[0]

			isbn := "000"
			if len(bookObj.Isbns) > 0 {
				isbn = bookObj.Isbns[0].Isbn10
			}

			if isbn == "000" || isbn == "" && len(bookObj.BookDetails) > 0 {
				isbn = bookObj.BookDetails[0].PrimaryIsbn10
			}

			week := bookObj.RankLastWeek
			books = append(books, SimpleBook{
				LastWeekRank:     week,
				WeeksOnList:      bookObj.WeeksOnList,
				Rank:             bookObj.Rank,
				AmazonProductURL: bookObj.AmazonProductURL,
				Publisher:        bookInfo.Publisher,
				Title:            bookInfo.Title,
				Author:           bookInfo.Author,
				Description:      bookInfo.Description,
				Isbn:             isbn,
			})
		}
	}

	return books
}

func (r *RestClient) GetBookCoverURL(isbn string) BookCoverURL {
	storedURLIf, _ := r.isbnCoverURLMap.Load(isbn)
	if storedURL, ok := storedURLIf.(string); ok {
		if storedURL != "" {
			return BookCoverURL{URL: storedURL}
		}
	}

	defaultURL := BookCoverURL{
		URL: "https://cos143.y3sh.com/bookPlaceholder.png",
	}
	url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=isbn:%s&key=%s", isbn, r.googleBookAPIKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Could not make req for googleBookRes, %s", err.Error())
		return defaultURL
	}

	res, err := r.httpClient.Do(req)
	if err != nil {
		log.Errorf("Could not fetch googleBookRes, %s", err.Error())
		return defaultURL
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("Non 200 status on fetch googleBookRes: , %d", res.StatusCode)
		return defaultURL
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Errorf("Could not read googleBookRes, %s", err.Error())
		return defaultURL
	}

	googleBookRes := GoogleBookRes{}
	jsonErr := json.Unmarshal(body, &googleBookRes)
	if jsonErr != nil {
		log.Errorf("Could not unmarshall googleBookRes, %s", err.Error())
		return defaultURL
	}

	items := googleBookRes.Items
	thumbnail := defaultURL.URL
	if len(items) > 0 {
		thumbnail = items[0].VolumeInfo.ImageLinks.Thumbnail
	}

	r.isbnCoverURLMap.Store(isbn, thumbnail)

	return BookCoverURL{URL: thumbnail}
}

func (r *RestClient) getBestSellers() BestSellerRes {
	if len(r.lastBestSellers.Results) > 0 {
		bTimeStr := r.lastBestSellers.Results[0].BestsellersDate
		bTime, _ := time.Parse("2006-01-02", bTimeStr)
		now := time.Now()
		if now.Sub(bTime).Hours() > 6 {
			return r.lastBestSellers
		}
	}

	url := fmt.Sprintf("https://api.nytimes.com/svc/books/v3/lists.json?list-name=hardcover-fiction&api-key=%s", r.bestSellerAPIKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Could not make req for best sellers, %s", err.Error())
		return r.lastBestSellers
	}

	res, err := r.httpClient.Do(req)
	if err != nil {
		log.Errorf("Could not fetch best sellers, %s", err.Error())
		return r.lastBestSellers
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("Non 200 status on fetch best sellers: , %d", res.StatusCode)
		return r.lastBestSellers
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Errorf("Could not read best sellers, %s", err.Error())
		return r.lastBestSellers
	}

	bestSellerRes := BestSellerRes{}
	jsonErr := json.Unmarshal(body, &bestSellerRes)
	if jsonErr != nil {
		log.Errorf("Could not unmarshall best sellers, %s", err.Error())
		return r.lastBestSellers
	}

	r.lastBestSellers = bestSellerRes

	return bestSellerRes
}
