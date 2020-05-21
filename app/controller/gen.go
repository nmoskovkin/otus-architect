package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/repository"
	"architectSocial/domain"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

func CreateGeneratorGetHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	mutex := sync.Mutex{}

	return func(w http.ResponseWriter, r *http.Request) error {
		mutex.Lock()
		fmt.Fprintln(w, "queued")
		dg := helpers.CreateRandomUserGenerator("u", 0)
		channel := make(chan []domain.RegisterUserDto)
		doneChannel := make(chan bool)
		registerService := domain.CreateRegisterManyUsersService(repository.CreateMysqlUserRepository(db))
		counterMutex := sync.Mutex{}
		counter := 0

		addCounter := func(v int) {
			counterMutex.Lock()
			counter += v
			if counter%100 == 0 {
				fmt.Printf("Counter: %d\n", counter)
			}

			counterMutex.Unlock()
		}

		worker := func(channel chan []domain.RegisterUserDto, doneChannel chan bool) {
			for v := range channel {
				validationResult, err := registerService(&domain.RegisterManyUsersDto{Users: v})
				if validationResult != nil && !validationResult.IsValid() {
					fmt.Printf("generate user error: validation failed %v\n", validationResult.GetAllErrors())
					continue
				}
				if err != nil {
					fmt.Printf("generate user error: %s\n", err.Error())
					continue
				}

				addCounter(len(v))
			}
			doneChannel <- true
		}

		go dg.Generate(channel, 1000, 4000)
		workerCount := 4

		for i := 0; i < workerCount; i++ {
			go worker(channel, doneChannel)
		}
		go func() {
			for i := 0; i < workerCount; i++ {
				<-doneChannel
			}
			mutex.Unlock()
		}()

		return nil
	}
}
