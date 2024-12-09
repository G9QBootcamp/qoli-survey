package dto

import "github.com/G9QBootcamp/qoli-survey/internal/user/models"

type SurveyRoleAssignRequest struct {
	SurveyID      uint   `json:"survey_id"`
	UserID        uint   `json:"user_id"`
	RoleName      string `json:"role_name"`
	PermissionIds []uint `json:"permission_ids" validate:"required"`
	TimeLimit     *int   `json:"time_limit,omitempty" validate:"omitempty,min=0"`
}

type SurveyRoleAssignResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	SurveyID    uint                `json:"survey_id"`
	RoleID      uint                `json:"role_id"`
	Permissions []models.Permission `json:"permissions"`
	TimeLimit   int                 `json:"time_limit"`
}
type GetUserRolesForSomeSurveyResponse struct {
	UserID   uint   `json:"user_id"`
	SurveyID uint   `json:"survey_id"`
	Roles    []Role `json:"roles"`
}

type Role struct {
	ID          uint         `json:"id"`
	Permissions []Permission `json:"permissions"`
	TimeLimit   int          `json:"time_limit"`
}

type Permission struct {
	Action string `json:"action"`
}
