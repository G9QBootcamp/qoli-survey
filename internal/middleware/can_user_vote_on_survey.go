package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	surveyModels "github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	userModels "github.com/G9QBootcamp/qoli-survey/internal/user/models"

	"github.com/labstack/echo/v4"
)

func CanUserVoteOnSurvey(db db.DbService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userID, ok := c.Get("userID").(uint)
			if !ok || userID == 0 {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "user not found"})
			}

			var user userModels.User

			if err := db.GetDb().First(&user, "id = ?", userID).Error; err != nil {
				return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "user not found"})
			}

			surveyID := c.Param("survey_id")
			if surveyID == "" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "survey id required"})
			}
			var survey surveyModels.Survey
			if err := db.GetDb().Preload("Options").First(&survey, "id = ?", surveyID).Error; err != nil {
				return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "survey not found"})
			}

			fmt.Println(survey.Options)
			for _, v := range survey.Options {
				if v.Name == "city" && v.Value != user.City {
					return echo.NewHTTPError(http.StatusForbidden, map[string]string{"error": "you are not allowed to participate to survey because of your city"})
				}

				if v.Name == "minimum_age" {
					require_age, err := strconv.Atoi(v.Value)
					if err != nil {
						continue
					}

					age := time.Now().Year() - user.DateOfBirth.Year()

					if age < require_age {
						return echo.NewHTTPError(http.StatusForbidden, map[string]string{"error": "you are not allowed to participate to survey because og your age"})

					}

				}
			}
			return next(c)

		}
	}
}
