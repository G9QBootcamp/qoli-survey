package middlewares

import (
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	surveyModels "github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	userModels "github.com/G9QBootcamp/qoli-survey/internal/user/models"

	"github.com/labstack/echo/v4"
)

func CheckPermission(requiredPermission string, db db.DbService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userID, ok := c.Get("userID").(uint)
			if !ok || userID == 0 {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
			}

			var user userModels.User
			if err := db.GetDb().Preload("GlobalRole.Permissions").First(&user, userID).Error; err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"error": "user not found"})
			}

			for _, permission := range user.GlobalRole.Permissions {
				if permission.Action == requiredPermission {
					return next(c)
				}
			}

			surveyID := c.Param("survey_id")
			if surveyID != "" {
				var survey surveyModels.Survey
				if err := db.GetDb().First(&survey, "id = ?", surveyID).Error; err != nil {
					return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "survey not found"})
				}

				if survey.OwnerID == userID {
					return next(c)
				}

				var userSurveyRoles []userModels.UserSurveyRole
				if err := db.GetDb().Preload("Role.Permissions").
					Where("user_id = ? AND survey_id = ?", userID, surveyID).
					Find(&userSurveyRoles).Error; err != nil {
					return echo.NewHTTPError(http.StatusForbidden, map[string]string{"error": "access denied for this survey"})
				}

				for _, userSurveyRole := range userSurveyRoles {
					for _, permission := range userSurveyRole.Role.Permissions {
						if permission.Action == requiredPermission {
							return next(c)
						}
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, map[string]string{"error": "access denied"})
		}
	}
}
