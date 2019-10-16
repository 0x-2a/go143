package twitter

import (
	"github.com/juju/errors"
)

type Tweet struct {
	TweetText string `json:"tweetText"`
}

type TweetService struct {
	Tweets []*Tweet `json:"tweets"`
}

func NewTweetService() *TweetService {
	return &TweetService{
		Tweets: []*Tweet{},
	}
}

func (t *TweetService) GetTweets() []*Tweet {
	return t.Tweets
}

func (t *TweetService) AddTweet(tweet *Tweet) error {
	if tweet == nil {
		return errors.New("missing tweet")
	}

	// Prepend the tweet
	t.Tweets = append([]*Tweet{tweet}, t.Tweets...)

	return nil
}
