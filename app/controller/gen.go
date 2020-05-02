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

func debug(dto domain.RegisterUserDto) {
	fmt.Println("----")
	fmt.Printf("firstName: %s\n", dto.FirstName)
	fmt.Printf("lastName: %s\n", dto.LastName)
	fmt.Printf("age: %s\n", dto.Age)
	fmt.Printf("city: %s\n", dto.City)
	fmt.Printf("login: %s\n", dto.Login)
	fmt.Printf("gender: %s\n", dto.Gender)
	fmt.Printf("password: %s\n", dto.Password)
	fmt.Printf("passwordConf: %s\n", dto.PasswordConfirmation)
	fmt.Printf("interests: %s\n", dto.Interests)

	fmt.Println("----")
}

func debugMany(dtos []domain.RegisterUserDto) {
	for _, r := range dtos {
		debug(r)
	}
}

func createGenFunc(in chan []domain.RegisterUserDto, registerManyUsersService domain.RegisterManyUsersService) {
	f := func(in chan []domain.RegisterUserDto) {
		for users := range in {
			//debugMany(users)
			validationResult, err := registerManyUsersService(&domain.RegisterManyUsersDto{
				Users: users,
			})
			//if !validationResult.IsValid() {
			//	panic(fmt.Sprintln(validationResult))
			//}
			if err != nil {
				fmt.Println(err)
				//w.Write([]byte(err.Error()))
			}
			if err != nil && validationResult != nil && !validationResult.IsValid() {
				fmt.Println(err)
				//w.Write([]byte("Invalid"))
			}
		}
	}

	go f(in)
}

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
		in := make(chan []domain.RegisterUserDto)
		defer close(in)

		createGenFunc(in, registerService)
		createGenFunc(in, registerService)
		//createGenFunc(in, registerService)
		//createGenFunc(in, registerService)

		c := 0
		// Incorrect behaviour
		// Data are generated but sql queries aren't waited for finish.
		generator := helpers.CreateRandomUserGenerator(func(i int) {
			_, _ = w.Write([]byte(fmt.Sprintf("Completed: %d", i)))
		}, func(users []domain.RegisterUserDto) {
			c += len(users)
			fmt.Printf("C: %d\n\n\n", c)
			//debugMany(users)
			in <- users
		})
		err = generator(1000000)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}
