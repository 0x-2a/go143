package instagram

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
)

type User struct {
	MobileEmail string   `json:"mobileEmail"`
	FullName    string   `json:"fullName"`
	Username    string `json:"username"`
	Password    string   `json:"password"`
}

type RandomUser struct {
	Name       string   `json:"name"`
	Location   string   `json:"location"`
	Picture    string   `json:"picture"`
	FeedImages []string `json:"feedImages"`
}

type UserKey struct {
	cseName string
	userName string
}

type UserService struct {
	userMutex *sync.Mutex
	UserMap   map[UserKey]*User
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
		UserMap:   make(map[UserKey]*User),
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

func (u *UserService) AddUser(cseName string, user User) error {
	if strings.TrimSpace(user.MobileEmail) == "" ||
		strings.TrimSpace(user.Username) == "" ||
		strings.TrimSpace(user.Password) == "" {
		return errors.New("missing required user fields")
	}

	u.userMutex.Lock()
	defer u.userMutex.Unlock()

	key := UserKey{
		cseName:  cseName,
		userName: user.Username,
	}

	if _, ok := u.UserMap[key]; ok {
		return errors.New("user already exists")
	}

	u.UserMap[key] = &user

	return nil
}

func (u *UserService) GetUsers(cseName string) []User {
	users := []User{}

	for k := range u.UserMap {
		if k.cseName == cseName{
			users = append(users, *u.UserMap[k])
		}
	}

	return users
}

func (u *UserService) IsValidPassword(cseName, username, passwordAttempt string) bool {
	u.userMutex.Lock()
	defer u.userMutex.Unlock()

	key := UserKey{
		cseName:  cseName,
		userName: username,
	}
	if user, ok := u.UserMap[key]; ok {
		return user.Password == passwordAttempt
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
