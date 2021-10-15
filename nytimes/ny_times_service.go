package nytimes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

var (
	httpClient = http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
)

type RestClient struct {
	apiKey          string
	lastBestSellers BestSellerRes
}

// 6ad84e249d054efeaefe1abb8f89df5b
func NewRestClient(apiKey string) *RestClient {
	return &RestClient{
		apiKey: apiKey,
	}
}

func (r *RestClient) GetBestSellers() BestSellerRes {
	if len(r.lastBestSellers.Results) > 0 {
		bTimeStr := r.lastBestSellers.Results[0].BestsellersDate
		bTime, _ := time.Parse("2006-01-02", bTimeStr)
		now := time.Now()
		if now.Sub(bTime).Hours() > 6 {
			return r.lastBestSellers
		}
	}

	url := fmt.Sprintf("https://api.nytimes.com/svc/books/v3/lists.json?list-name=hardcover-fiction&api-key=%s", r.apiKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Could not make req for best sellers, %s", err.Error())
		return r.lastBestSellers
	}

	res, err := httpClient.Do(req)
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
