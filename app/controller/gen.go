package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/repository"
	"architectSocial/domain"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func CreateGeneratorGetHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		isAuth, _, err := sessionWrapper.IsAuthenticated(r)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		if !isAuth {
			return NewHTTPError(fmt.Errorf("forbidden"), 403, "")
		}
		registerService := domain.CreateRegisterManyUsersService(repository.CreateMysqlUserRepository(db))
		generator := helpers.CreateRandomUserGenerator(func(i int) {
			_, _ = w.Write([]byte(fmt.Sprintf("Completed: %d", i)))
		}, func(users []domain.RegisterUserDto) {
			fmt.Println(users)
			validationResult, err := registerService(&domain.RegisterManyUsersDto{
				Users: users,
			})
			//if !validationResult.IsValid() {
			//	panic(fmt.Sprintln(validationResult))
			//}
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			if err != nil && validationResult != nil && !validationResult.IsValid() {
				w.Write([]byte("Invalid"))
			}
		})
		err = generator(1)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}
