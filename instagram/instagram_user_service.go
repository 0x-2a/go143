package instagram

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"strings"
	"sync"
)

type User struct {
	MobileEmail string   `json:"mobileEmail"`
	FullName    string   `json:"fullName"`
	Username    Username `json:"username"`
	Password    string   `json:"password"`
}

type RandomUser struct {
	Name       string   `json:"name"`
	Location   string   `json:"location"`
	Picture    string   `json:"picture"`
	FeedImages []string `json:"feedImages"`
}

type Username string

type UserService struct {
	userMutex *sync.Mutex
	UserMap   map[Username]*User
}

var (
	randUsers = []RandomUser{
		{
			Name:     "Miss Addison Young",
			Location: "6682 Brock Rd, Fountainbleu, Nunavut, Canada, S8P 4T7",
			Picture:  "https://cos143.y3sh.com/user1.jpg",
		},
		{
			Name:     "Miss Alice Spencer",
			Location: "922 Frostfield Dr, New York, NY, USA",
			Picture:  "https://cos143.y3sh.com/user2.jpg",
		},
		{
			Name:     "Miss Vallery Kirkbride",
			Location: "8472 Connifer Ridge Rd, Mountain View, CA",
			Picture:  "https://cos143.y3sh.com/user3.jpg",
		},
		{
			Name:     "Mr Jordan Montoya",
			Location: "369 Roam Terrace, Atlanta, GA, USA",
			Picture:  "https://cos143.y3sh.com/user4.jpg",
		},
		{
			Name:     "Mr Lucas Bryant",
			Location: "6395 Wheathill Pass, Boulder, CO, USA",
			Picture:  "https://cos143.y3sh.com/user5.jpg",
		},
		{
			Name:     "Mr Frank Anderson",
			Location: "455 Benton Blvd, San Francisco, CA, USA",
			Picture:  "https://cos143.y3sh.com/user6.jpg",
		},
		{
			Name:     "Mr Bernard Abernathy",
			Location: "221B Easy St, Mountain View, CA, USA",
			Picture:  "https://cos143.y3sh.com/user7.jpg",
		},
	}

	feedImages = []string{
		"https://cos143.y3sh.com/insta1.jpg",
		"https://cos143.y3sh.com/insta2.jpg",
		"https://cos143.y3sh.com/insta3.jpg",
		"https://cos143.y3sh.com/insta4.jpg",
		"https://cos143.y3sh.com/insta5.jpg",
		"https://cos143.y3sh.com/insta6.jpg",
		"https://cos143.y3sh.com/insta7.jpg",
		"https://cos143.y3sh.com/insta8.jpg",
		"https://cos143.y3sh.com/insta9.jpg",
		"https://cos143.y3sh.com/insta10.jpg",
		"https://cos143.y3sh.com/insta11.jpg",
		"https://cos143.y3sh.com/insta12.jpg",
		"https://cos143.y3sh.com/insta13.jpg",
		"https://cos143.y3sh.com/insta14.jpg",
		"https://cos143.y3sh.com/insta15.jpg",
		"https://cos143.y3sh.com/insta16.jpg",
		"https://cos143.y3sh.com/insta17.jpg",
		"https://cos143.y3sh.com/insta18.jpg",
		"https://cos143.y3sh.com/insta19.jpg",
	}
)

func NewUserService() *UserService {
	return &UserService{
		UserMap:   make(map[Username]*User),
		userMutex: &sync.Mutex{},
	}
}

func (u *UserService) GetTweets() []*User {
	u.userMutex.Lock()
	defer u.userMutex.Unlock()

	var users []*User
	for _, v := range u.UserMap {
		users = append(users, v)
	}

	return users
}

func (u *UserService) AddUser(user User) error {
	if strings.TrimSpace(user.MobileEmail) == "" ||
		strings.TrimSpace(string(user.Username)) == "" ||
		strings.TrimSpace(user.Password) == "" {
		return errors.New("missing required user fields")
	}

	lowSecurityHash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(lowSecurityHash[:])

	u.userMutex.Lock()
	defer u.userMutex.Unlock()

	if _, ok := u.UserMap[user.Username]; ok {
		return errors.New("user already exists")
	}

	u.UserMap[user.Username] = &user

	return nil
}

func (u *UserService) GetUsers() []User {
	users := []User{}

	for k := range u.UserMap {
		user := *u.UserMap[k]

		users = append(users, user)
	}

	return users
}

func (u *UserService) IsValidPassword(username Username, passwordAttempt string) bool {
	lowSecurityHash := md5.Sum([]byte(passwordAttempt))
	passwordAttemptHash := hex.EncodeToString(lowSecurityHash[:])

	u.userMutex.Lock()
	defer u.userMutex.Unlock()

	if user, ok := u.UserMap[username]; ok {
		return user.Password == passwordAttemptHash
	}

	return false
}

func (u *UserService) GetRandProfile() RandomUser {
	user := randUsers[rand.Intn(len(randUsers))]

	p := rand.Perm(9)

	for _, r := range p {
		user.FeedImages = append(user.FeedImages, feedImages[r])
	}

	return user
}
