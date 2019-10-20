package instagram

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"sync"

	"github.com/juju/errors"
)

type User struct {
	MobileEmail string   `json:"mobileEmail"`
	FullName    string   `json:"fullName"`
	Username    Username `json:"username"`
	Password    string   `json:"password"`
}

type Username string

type UserService struct {
	userMutex *sync.Mutex
	UserMap   map[Username]*User
}

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
