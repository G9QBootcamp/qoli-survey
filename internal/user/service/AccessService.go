package service

import (
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	notification "github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type IAccessService interface {
	SetRole(c context.Context, req dto.SurveyRoleAssignRequest) (*dto.SurveyRoleAssignResponse, error)
	GetUserRolesForSomeSurvey(c context.Context, userID uint, surveyID uint) (*dto.GetUserRolesForSomeSurveyResponse, error)
	GetAllPermissions(c context.Context) ([]models.Permission, error)
	DeleteUserSurveyRole(c context.Context, surveyID uint, userID uint, roleID uint) error
}
type AccessService struct {
	conf                *config.Config
	repo                repository.IAccessRepository
	logger              logging.Logger
	notificationService notification.INotificationService
}

func NewAccessService(conf *config.Config, repo repository.IAccessRepository, logger logging.Logger, notificationService notification.INotificationService) *AccessService {
	return &AccessService{conf: conf, repo: repo, logger: logger, notificationService: notificationService}
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
		Name:        req.RoleName,
		Permissions: prms,
	}
	createdRole, err := s.repo.CreateRole(c, role)
	if err != nil {
		return nil, err
	}
	// usr is the abbreviation for UserSurveyRole
	usr := models.UserSurveyRole{
		UserID:   req.UserID,
		SurveyID: req.SurveyID,
		RoleID:   createdRole.ID,
	}
	if req.TimeLimit != nil {
		usr.TimeLimit = *req.TimeLimit
	}
	createdUsr, err := s.repo.CreateUserSurveyRole(c, usr)
	if err != nil {
		return nil, err
	}
	_, err = s.notificationService.Notify(c, usr.UserID, "new role : "+createdRole.Name+" assigned to you.")
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToSendNotify, "error in sending notify in access service", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
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

func (s *AccessService) GetUserRolesForSomeSurvey(c context.Context, userID uint, surveyID uint) (*dto.GetUserRolesForSomeSurveyResponse, error) {
	res, err := s.repo.GetUserRolesForSurvey(c, userID, surveyID)
	if err != nil {
		return nil, err
	}
	response := dto.GetUserRolesForSomeSurveyResponse{
		UserID:   userID,
		SurveyID: surveyID,
	}

	for _, role := range res {
		var permissionsDTO []dto.Permission
		for _, permission := range role.Role.Permissions {
			permissionsDTO = append(permissionsDTO, dto.Permission{
				Action: permission.Action,
			})
		}

		response.Roles = append(response.Roles, dto.Role{
			ID:          role.RoleID,
			Permissions: permissionsDTO,
			TimeLimit:   role.TimeLimit,
		})
	}
	return &response, nil
}
func (s *AccessService) GetAllPermissions(c context.Context) ([]models.Permission, error) {
	return s.repo.GetAllPermissions(c)
}
func (s *AccessService) DeleteUserSurveyRole(c context.Context, surveyID uint, userID uint, roleID uint) error {
	return s.repo.DeleteUserSurveyRole(c, surveyID, userID, roleID)
}
