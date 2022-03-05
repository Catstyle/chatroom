package services

import (
	"io/ioutil"
	"strings"

	"github.com/catstyle/chatroom/utils"
)

type sensitiveService struct {
	filter *utils.TrieTree
}

var sensitiveSvc *sensitiveService

func InitSensitiveService() {
	options := utils.GetOptions()
	tt := utils.NewTrieTree()

	for _, filename := range options.GetStringSlice("PROFANITY_TEXT") {
		if text, err := ioutil.ReadFile(filename); err == nil {
			tt.AddWord(strings.Split(string(text), "\n")...)
		}
	}

	sensitiveSvc = &sensitiveService{
		filter: tt,
	}
}

func GetSensitiveService() *sensitiveService {
	return sensitiveSvc
}

func (svc *sensitiveService) Filter(text string, mask rune) string {
	return svc.filter.Filter(text, mask)
}
