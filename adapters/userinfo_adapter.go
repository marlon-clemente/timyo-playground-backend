package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/ports"
)

type userInfo struct {
	authServiceURL string
}

func NewUserInfo(authServiceURL string) ports.UserInfoService {
	return &userInfo{
		authServiceURL: authServiceURL,
	}
}

func (s *userInfo) GetUserInfo(ctx context.Context, memberID string, accessToken string) (*ports.UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.authServiceURL+"/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build user info request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 8 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call user info service: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("user info service returned status %d", res.StatusCode)
	}

	type userInfoResponse struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		Avatar    string `json:"avatar"`
	}

	var payload userInfoResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %w", err)
	}

	resolvedID := strings.TrimSpace(payload.ID)
	resolvedName := strings.TrimSpace(payload.Name)
	resolvedAvatar := strings.TrimSpace(payload.Avatar)

	return &ports.UserInfo{
		ID: resolvedID,
		Name: resolvedName,
		AvatarURL: resolvedAvatar,
	}, nil
}
