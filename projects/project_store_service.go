package projects

import (
	"fmt"
	"sync"
)

type User struct {
	MobileEmail string   `json:"mobileEmail"`
	FullName    string   `json:"fullName"`
	Username    Username `json:"username"`
	Password    string   `json:"password"`
}

type Username string

type ProjectStoreService struct {
	storeMutex   *sync.Mutex
	ProjectStore map[string]string
}

func NewProjectStoreService() *ProjectStoreService {
	return &ProjectStoreService{
		ProjectStore: make(map[string]string),
		storeMutex:   &sync.Mutex{},
	}
}

func (u *ProjectStoreService) GetValue(groupName, keyName string) string {
	u.storeMutex.Lock()
	defer u.storeMutex.Unlock()

	val, ok := u.ProjectStore[fmt.Sprintf("%s:%s", groupName, keyName)]
	if ok {
		return val
	}

	return ""
}

func (u *ProjectStoreService) SetValue(groupName, keyName, value string) {
	u.storeMutex.Lock()
	defer u.storeMutex.Unlock()

	u.ProjectStore[fmt.Sprintf("%s:%s", groupName, keyName)] = value
}
