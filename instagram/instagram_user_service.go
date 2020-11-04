package instagram

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
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
	Name     string `json:"name"`
	Location string `json:"location"`
	Picture  string `json:"picture"`
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
			Picture:  "https://randomuser.me/api/portraits/women/3.jpg",
		},
		{
			Name:     "Miss Addison Young",
			Location: "6682 Brock Rd, Fountainbleu, Nunavut, Canada, S8P 4T7",
			Picture:  "https://randomuser.me/api/portraits/women/3.jpg",
		},
		{
			Name:     "Miss Addison Young",
			Location: "6682 Brock Rd, Fountainbleu, Nunavut, Canada, S8P 4T7",
			Picture:  "https://randomuser.me/api/portraits/women/3.jpg",
		},
		{
			Name:     "Miss Addison Young",
			Location: "6682 Brock Rd, Fountainbleu, Nunavut, Canada, S8P 4T7",
			Picture:  "https://randomuser.me/api/portraits/women/3.jpg",
		},
		{
			Name:     "Miss Addison Young",
			Location: "6682 Brock Rd, Fountainbleu, Nunavut, Canada, S8P 4T7",
			Picture:  "https://randomuser.me/api/portraits/women/3.jpg",
		},
		{
			Name:     "Miss Addison Young",
			Location: "6682 Brock Rd, Fountainbleu, Nunavut, Canada, S8P 4T7",
			Picture:  "https://randomuser.me/api/portraits/women/3.jpg",
		},
		{
			Name:     "Miss Addison Young",
			Location: "6682 Brock Rd, Fountainbleu, Nunavut, Canada, S8P 4T7",
			Picture:  "https://randomuser.me/api/portraits/women/3.jpg",
		},
	}

	// 3 girls, 4 boys

	userImages = []string{
		"https://cos143.s3.us-east-2.amazonaws.com/user1.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/user2.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/user3.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/user4.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/user5.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/user6.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/user7.jpg",
	}

	feedImages = []string{
		"https://cos143.s3.us-east-2.amazonaws.com/insta1.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta2.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta3.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta4.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta5.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta6.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta7.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta8.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta9.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta10.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta11.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta12.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta13.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta14.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta15.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta16.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta17.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta18.jpg",
		"https://cos143.s3.us-east-2.amazonaws.com/insta19.jpg",
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

func (u *UserService) GetRandProfile() {

}
