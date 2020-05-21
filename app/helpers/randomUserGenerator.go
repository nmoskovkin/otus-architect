package helpers

import (
	"architectSocial/domain"
)

type RandomUserGenerator struct {
	counter      int
	dataProvider *RandomUserDataProvider
}

func CreateRandomUserGenerator(prefix string, counter int) *RandomUserGenerator {
	dp, _ := CreateRandomUserDataProvider(prefix)
	return &RandomUserGenerator{
		counter:      counter,
		dataProvider: dp,
	}
}

func (dp *RandomUserGenerator) Generate(channel chan []domain.RegisterUserDto, bunchSize int, bunchCount int) {
	counter := dp.counter
	for i := 0; i < bunchCount; i++ {
		var dtos []domain.RegisterUserDto
		for j := 0; j < bunchSize; j++ {
			dtos = append(dtos, domain.RegisterUserDto{
				FirstName:            dp.dataProvider.FirstName(),
				LastName:             dp.dataProvider.LastName(),
				Age:                  dp.dataProvider.Age(),
				Gender:               dp.dataProvider.Gender(),
				Interests:            dp.dataProvider.Text(),
				City:                 dp.dataProvider.City(),
				Password:             dp.dataProvider.Login(counter),
				PasswordConfirmation: dp.dataProvider.Login(counter),
				Login:                dp.dataProvider.Login(counter),
			})
			counter++
		}
		channel <- dtos
	}

	close(channel)
}
