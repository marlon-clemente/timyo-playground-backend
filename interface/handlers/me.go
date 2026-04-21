package handlers

import (
	"github.com/marlon-clemente/timyo-playground-backend/internal/application/queries"
	"github.com/marlon-clemente/timyo-playground-backend/packages/server"
)

type MeHandler struct {
	queryMe queries.Me
}

func NewMeHandler(q queries.Me) *MeHandler {
	return &MeHandler{queryMe: q}
}

func (h *MeHandler) GetMe(c *server.Ctx) error {
	memberId := c.GetUserID()
	q, err := h.queryMe.Me(c.Context(), memberId)
	if err != nil {
		return err
	}
	return c.ResponseOk(q)
}