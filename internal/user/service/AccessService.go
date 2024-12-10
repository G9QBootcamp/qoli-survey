package service

import (
	"errors"
	"time"

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
	CreateVoteVisibility(c context.Context, surveyID uint, viewerID uint, request dto.VoteVisibilityCreateRequest) (dto.VoteVisibilityResponse, error)
	//GetVoteVisibilityById(ctx context.Context, id uint) (dto.VoteVisibilityResponse, error)
	//GetVoteVisibilityBySurveyId(ctx context.Context, surveyId uint) ([]dto.VoteVisibilityResponse, error)
	//DeleteVoteVisibilityById(ctx context.Context, id uint) error
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
		usr.ExpiresAt = time.Now().Add(time.Duration(*req.TimeLimit) * time.Minute)
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
		ExpiresAt:   createdUsr.ExpiresAt,
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

		roleModel, err := s.repo.GetRoleByID(c, role.RoleID)
		if err != nil {
			return nil, err
		}

		response.Roles = append(response.Roles, dto.Role{
			ID:          roleModel.ID,
			Name:        roleModel.Name,
			Permissions: permissionsDTO,
			ExpiresAt:   role.ExpiresAt,
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
func (s *AccessService) CreateVoteVisibility(c context.Context, surveyID uint, viewerID uint, req dto.VoteVisibilityCreateRequest) (dto.VoteVisibilityResponse, error) {
	for _, respondentID := range req.RespondentIDs {
		_, err := s.repo.CreateVoteVisibility(c, surveyID, viewerID, uint(respondentID))
		if err != nil {
			return dto.VoteVisibilityResponse{}, err
		}
	}

	return dto.VoteVisibilityResponse{
		SurveyID:      int(surveyID),
		ViewerID:      int(viewerID),
		RespondentIDs: req.RespondentIDs,
	}, nil
}

//	func (s *AccessService) GetVoteVisibilityById(c context.Context, id uint) (dto.VoteVisibilityResponse, error) {
//		res, err := s.repo.GetVoteVisibilityById(c, id)
//		if err != nil {
//			return dto.VoteVisibilityResponse{}, err
//		}
//		return dto.VoteVisibilityResponse{
//			ID:           res.ID,
//			SurveyID:     res.SurveyID,
//			ViewerID:     res.ViewerID,
//			RespondentID: res.RespondentID,
//		}, err
//	}
//
//	func (s *AccessService) GetVoteVisibilityBySurveyId(ctx context.Context, surveyId uint) ([]dto.VoteVisibilityResponse, error) {
//		res, err := s.repo.GetVoteVisibilityBySurveyId(ctx, surveyId)
//		if err != nil {
//			return nil, err
//		}
//		var response []dto.VoteVisibilityResponse
//		for _, vv := range res {
//			response = append(response, dto.VoteVisibilityResponse{
//				ID:           vv.ID,
//				SurveyID:     vv.SurveyID,
//				ViewerID:     vv.ViewerID,
//				RespondentID: vv.RespondentID,
//			})
//		}
//		return response, nil
//	}
func (s *AccessService) DeleteVoteVisibilityById(c context.Context, id uint) error {
	return s.repo.DeleteVoteVisibilityById(c, id)
}
