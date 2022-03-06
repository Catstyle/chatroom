package services

import "github.com/catstyle/chatroom/pkg/models"

type StatsService interface {
	OnlineStats(string) *models.OnlineUser
}

type statsService struct {
}

var statsSvc *statsService

func NewStatsService() {
}

func GetStatsService() StatsService {
	return statsSvc
}

func (svc *statsService) OnlineStats(username string) *models.OnlineUser {
	return nil
}
