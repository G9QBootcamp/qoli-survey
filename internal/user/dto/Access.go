package dto

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
)

type SurveyRoleAssignRequest struct {
	SurveyID      uint   `json:"survey_id"`
	UserID        uint   `json:"user_id"`
	RoleName      string `json:"role_name" validate:"required"`
	PermissionIds []uint `json:"permission_ids" validate:"required"`
	TimeLimit     *int   `json:"time_limit,omitempty" validate:"omitempty,min=0"`
}

type SurveyRoleAssignResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	SurveyID    uint                `json:"survey_id"`
	RoleID      uint                `json:"role_id"`
	Permissions []models.Permission `json:"permissions"`
	ExpiresAt   time.Time           `json:"expires_at"`
}
type GetUserRolesForSomeSurveyResponse struct {
	UserID   uint   `json:"user_id"`
	SurveyID uint   `json:"survey_id"`
	Roles    []Role `json:"roles"`
}

type Role struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	ExpiresAt   time.Time    `json:"expires_at"`
}

type Permission struct {
	Action string `json:"action"`
}

type VoteVisibilityCreateRequest struct {
	RespondentIDs []int `json:"respondent_ids" validate:"required"` // User whose votes can be viewed
}

type VoteVisibilityResponse struct {
	SurveyID      int   `json:"survey_id"`
	ViewerID      int   `json:"viewer_id"`      // User allowed to view votes
	RespondentIDs []int `json:"respondent_ids"` // User whose votes can be viewed
}
