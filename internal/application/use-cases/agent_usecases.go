package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/entities"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/ports"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/repositories"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
	"github.com/marlon-clemente/timyo-playground-backend/packages/errs"
)

type Agent struct {
	repositories repositories.Agents
	userInfoService ports.UserInfoService
}

func NewAgent(repositories repositories.Agents, userInfoService ports.UserInfoService) *Agent {
	return &Agent{
		repositories:    repositories,
		userInfoService: userInfoService,
	}
}

func (a *Agent) Create(ctx context.Context, input CreateAgentInput) (*AgentOutput, error) {
	memberId, err := vos.ParseUUID(input.MemberID)
	if err != nil {
		return nil, errs.DomainErr("Invalid member ID", err)
	}

	userInfo, err := a.userInfoService.GetUserInfo(ctx, input.MemberID, input.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	
	name, err := vos.NewName(userInfo.Name)
	if err != nil {
		return nil, errs.DomainErr("Invalid name", err)
	}

	ag := entities.NewAgent(memberId, name, userInfo.AvatarURL)
	slog.InfoContext(ctx, "agent created", "agent_id", ag.ID.String(), "member_id", ag.MemberID.String())

	err = a.repositories.Create(ctx, ag)
	if err != nil {
		return nil, errs.DomainErr("Failed to create agent", err)
	}

	return &AgentOutput{
		ID: ag.ID.String(),
		MemberID: ag.MemberID.String(),
		Name: ag.Name.String(),
		AvatarURL: ag.AvatarURL,
		CreatedAt: ag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}