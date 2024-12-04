package dto

import "github.com/G9QBootcamp/qoli-survey/internal/user/models"

type SurveyRoleAssignRequest struct {
	SurveyId      uint   `json:"surveyId" validate:"required"`
	UserId        uint   `json:"userId" validate:"required"`
	PermissionIds []uint `json:"permissionIds" validate:"required"`
	TimeLimit     *int   `json:"timeLimit" validate:"required"`
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
type GetUserRolesForSomeSurveyRequest struct {
	UserID   uint `json:"user_id" validate:"required"`
	SurveyID uint `json:"survey_id" validate:"required"`
}
type Role struct {
	Permissions []models.Permission `json:"permissions" validate:"required"`
	TimeLimit   int                 `json:"time_limit" validate:"required"`
}
type DeleteUserSurveyRoleRequest struct {
	ID uint `json:"id" validate:"required"`
}
