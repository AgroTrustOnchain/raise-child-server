package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type ILeaderNotiRepository interface {
	GetNoti(id string, ctx context.Context) (*entities.LeaderNoti, error)
	GetNotiByMealNeed(id string, ctx context.Context) (*entities.LeaderNoti, error)
	GetNotiByNeed(id string, ctx context.Context) (*entities.LeaderNoti, error)
	GetCurrentLeaderNotis(req request.GetNotisRequest, leader string, ctx context.Context) ([]entities.LeaderNoti, error)
	CreateNoti(noti entities.LeaderNoti, ctx context.Context) error
	UpdateNoti(noti entities.LeaderNoti, ctx context.Context) error
	AssignLeader(leader, region string, ctx context.Context) error
}
