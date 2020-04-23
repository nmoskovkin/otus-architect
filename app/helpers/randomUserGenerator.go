package helpers

import (
	"architectSocial/domain"
	"bufio"
	"encoding/csv"
	//"encoding/csv"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"strconv"
)

type RandomUserGenerator func(count int) error

const CountOfItemsToFlush = 2

type RandomGeneratorOutput func(i int)

type RandomSaver func(users []domain.RegisterUserDto)

type RandomGeneratorUser struct {
	firstName string
	lastName  string
	age       uint8
	city      string
	interests string
	gender    domain.UserGender
	id        uuid.UUID
	password  string
	salt      string
}

func CreateRandomUserGenerator(output RandomGeneratorOutput, saver RandomSaver) RandomUserGenerator {
	return func(count int) error {
		fpath := "data/first_names.all.txt"
		lpath := "data/last_names.all.txt"
		cpath := "data/world-cities-no-header.csv"
		tpath := "data/bbc-text-no-header.csv"

		f, errF := os.Open(fpath)
		l, errL := os.Open(lpath)
		c, errC := os.Open(cpath)
		t, errT := os.Open(tpath)
		fIter := 0
		lIter := 0
		cIter := 0
		tIter := 0
		bff := bufio.NewReader(f)
		bfl := bufio.NewReader(l)
		bfc := csv.NewReader(c)
		bft := csv.NewReader(t)

		if errF == nil {
			defer f.Close()
		}
		if errL == nil {
			defer l.Close()
		}
		if errC == nil {
			defer c.Close()
		}
		if errT == nil {
			defer t.Close()
		}
		getFirstName := func() string {
			v, _ := bff.ReadString('\n')
			fIter++
			if fIter == 164460 {
				f, _ = os.Open(fpath)
				bff = bufio.NewReader(f)
				fIter = 0
			}

			return v
		}
		getLastNameName := func() string {
			v, _ := bfl.ReadString('\n')
			lIter++
			if lIter == 98391 {
				l, _ = os.Open(lpath)
				bfl = bufio.NewReader(l)
				lIter = 0
			}

			return v
		}
		var getCity func() string
		getCity = func() string {
			v, _ := bfc.Read()
			cIter++
			if cIter == 23019 {
				c, _ = os.Open(cpath)
				bfc = csv.NewReader(c)
				cIter = 0
			}
			if len(v) == 0 {
				return getCity()
			}

			return v[0]
		}
		var getText func() string
		getText = func() string {
			v, _ := bft.Read()
			tIter++
			if tIter == 2225 {
				t, _ = os.Open(tpath)
				bft = csv.NewReader(t)
				tIter = 0
			}
			if len(v) == 0 {
				return getText()
			}
			return v[1]
		}

		//if errF != nil && errL != nil && errC != nil && errT != nil {
		//	return "", fmt.Printf("unable to read one of the files")
		//}
		//firstName := func(fn) {
		//	f, err := os.Open(fn)
		//	if err != nil {
		//		return "", err
		//	}
		//	defer f.Close()
		//}
		genderValues := []string{
			"male",
			"female",
		}
		randomUser := func(i int) domain.RegisterUserDto {
			//firstName, _ := ReadLineFromFile("data/first_names.all.txt", rand.Intn(164460))
			//lastName, _ := ReadLineFromFile("data/last_names.all.txt", rand.Intn(98391))
			//city, _ := ReadCsvCellFromFile("data/world-cities.csv", rand.Intn(23019), 0)
			//text, _ := ReadCsvCellFromFile("data/bbc-text.csv", rand.Intn(2225), 1)

			return domain.RegisterUserDto{
				FirstName:            getFirstName(),
				LastName:             getLastNameName(),
				Age:                  strconv.Itoa(rand.Intn(108)),
				Gender:               genderValues[rand.Intn(2)],
				Interests:            getText(),
				City:                 getCity(),
				Password:             "Test" + strconv.Itoa(i),
				PasswordConfirmation: "Test" + strconv.Itoa(i),
				Login:                "Test" + strconv.Itoa(i),
			}
		}

		flushNumber := 1
		lastIndex := 0
		users := []domain.RegisterUserDto{}
		for i := 1; i <= count; i++ {
			users = append(users, randomUser(i))
			if i%CountOfItemsToFlush == 0 {
				output(i)
				saver(users)
				flushNumber++
				users = []domain.RegisterUserDto{}
			}

			lastIndex = i
		}

		if lastIndex != flushNumber*CountOfItemsToFlush {
			output(lastIndex)
			saver(users)
			flushNumber++
		}

		return nil
	}
}
