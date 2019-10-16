package twitter

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/juju/errors"
)

const (
	tweetCapacity = 42
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Tweet struct {
	ID        int64  `json:"id"`
	TweetText string `json:"tweetText"`
	Timestamp int64  `json:"timestamp"`
}

type TweetService struct {
	tweetMutex *sync.Mutex
	Tweets     []*Tweet
	TweetCount int64
}

func NewTweetService() *TweetService {
	rand.Seed(time.Now().UnixNano())

	return &TweetService{
		Tweets:     []*Tweet{},
		tweetMutex: &sync.Mutex{},
	}
}

func (t *TweetService) GetTweets() []*Tweet {
	t.tweetMutex.Lock()
	defer t.tweetMutex.Unlock()

	return t.Tweets
}

func (t *TweetService) AddTweet(tweetText string) (*Tweet, error) {
	if strings.TrimSpace(tweetText) == "" {
		return nil, errors.New("missing tweet")
	}

	tweet := Tweet{
		TweetText: tweetText,
		Timestamp: time.Now().Unix(),
	}

	t.tweetMutex.Lock()
	defer t.tweetMutex.Unlock()

	t.TweetCount++
	tweet.ID = t.TweetCount

	t.Tweets = append(t.Tweets, &tweet)
	if len(t.Tweets) > tweetCapacity {
		t.Tweets = t.Tweets[1:]
	}

	return &tweet, nil
}

func (t *TweetService) AddRandTweet() (*Tweet, error) {
	return t.AddTweet(GetRandString(8))
}

func GetRandString(length int) string {
	runeBuf := make([]rune, length)

	for i := range runeBuf {
		runeBuf[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(runeBuf)
}
