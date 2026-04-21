package ports

import "context"

type UserInfo struct {
	ID string
	Name string
	AvatarURL string
}

type UserInfoService interface {
	GetUserInfo(ctx context.Context, memberID string, accessToken string) (*UserInfo, error)
}