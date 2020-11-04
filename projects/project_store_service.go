package projects

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type User struct {
	MobileEmail string   `json:"mobileEmail"`
	FullName    string   `json:"fullName"`
	Username    Username `json:"username"`
	Password    string   `json:"password"`
}

type Username string

type keyValRepository interface {
	SetKeyValue(key, value string) error
	GetValue(key string) (string, error)
}

type ProjectStoreService struct {
	keyValRepo keyValRepository
}

func NewProjectStoreService(keyValRepo keyValRepository) *ProjectStoreService {
	return &ProjectStoreService{
		keyValRepo: keyValRepo,
	}
}

func (p *ProjectStoreService) GetValue(groupName, keyName string) string {
	key := fmt.Sprintf("%s:%s", groupName, keyName)

	value, err := p.keyValRepo.GetValue(key)
	if err != nil {
		log.Errorf("Could not fetch key: %s\n%+v\n", key, err)
	}

	return value
}

func (p *ProjectStoreService) SetValue(groupName, keyName, value string) {

	key := fmt.Sprintf("%s:%s", groupName, keyName)

	err := p.keyValRepo.SetKeyValue(key, value)
	if err != nil {
		log.Errorf("Could not set key value: %s:%s\n%+v\n", key, value, err)
	}
}
