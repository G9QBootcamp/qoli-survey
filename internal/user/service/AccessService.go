package service

import (
	"errors"
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type IAccessService interface {
	SetRole(c context.Context, req dto.SurveyRoleAssignRequest) (*dto.SurveyRoleAssignResponse, error)
	GetUserRolesForSomeSurvey(c context.Context, userId, surveyId uint) (*dto.GetUserRolesForSomeSurveyResponse, error)
	GetAllPermissions(c context.Context) ([]models.Permission, error)
	DeleteUserSurveyRole(c context.Context, userId uint) error
}
type AccessService struct {
	conf   *config.Config
	repo   repository.IAccessRepository
	logger logging.Logger
}

func NewAccessService(conf *config.Config, repo repository.IAccessRepository, logger logging.Logger) *AccessService {
	return &AccessService{conf: conf, repo: repo, logger: logger}
}

func (s *AccessService) SetRole(c context.Context, req dto.SurveyRoleAssignRequest) (*dto.SurveyRoleAssignResponse, error) {
	allperms, err := s.repo.GetAllPermissions(c)
	if err != nil {
		return nil, err
	}
	var prms []models.Permission
	for _, pid := range req.PermissionIds {
		flag := false
		for _, prm := range allperms {
			if prm.ID == pid {
				flag = true
				prms = append(prms, prm)
				break
			}
		}
		if !flag {
			return nil, errors.New("permission not found")
		}
	}
	role := models.Role{
		Name:        "",
		Permissions: prms,
	}
	createdRole, err := s.repo.CreateRole(c, role)
	if err != nil {
		return nil, err
	}
	// usr is the abbreviation for UserSurveyRole
	usr := models.UserSurveyRole{
		UserID:   req.UserId,
		SurveyID: req.SurveyId,
		RoleID:   createdRole.ID,
	}
	if req.TimeLimit != nil {
		usr.TimeLimit = *req.TimeLimit
	}
	createdUsr, err := s.repo.CreateUserSurveyRole(c, usr)
	if err != nil {
		return nil, err
	}
	return &dto.SurveyRoleAssignResponse{
		ID:          createdUsr.ID,
		UserID:      createdUsr.UserID,
		SurveyID:    createdUsr.SurveyID,
		RoleID:      createdUsr.RoleID,
		Permissions: createdRole.Permissions,
		TimeLimit:   createdUsr.TimeLimit,
	}, nil
}

func (s *AccessService) GetUserRolesForSomeSurvey(c context.Context, userId, surveyId uint) (*dto.GetUserRolesForSomeSurveyResponse, error) {
	res, err := s.repo.GetUserRolesForSurvey(c, userId, surveyId)
	if err != nil {
		return nil, err
	}
	response := dto.GetUserRolesForSomeSurveyResponse{
		UserID:   userId,
		SurveyID: surveyId,
	}
	for _, role := range res {
		response.Roles = append(response.Roles, dto.Role{
			Permissions: role.Role.Permissions,
			TimeLimit:   role.TimeLimit,
		})
	}
	return &response, nil
}
func (s *AccessService) GetAllPermissions(c context.Context) ([]models.Permission, error) {
	return s.repo.GetAllPermissions(c)
}
func (s *AccessService) DeleteUserSurveyRole(c context.Context, userId uint) error {
	return s.repo.DeleteUserSurveyRole(c, userId)
}
