package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IVolunteerNotiRepository interface {
	GetNoti(id string, ctx context.Context) (*entities.VolunteerNoti, error)
	GetCurrentVolunteerNotis(req request.GetNotisRequest, volunteer string, ctx context.Context) ([]entities.VolunteerNoti, error)
	CreateNoti(noti entities.VolunteerNoti, ctx context.Context) error
	AssignVolunteer(volunteer, region string, ctx context.Context) error
}
