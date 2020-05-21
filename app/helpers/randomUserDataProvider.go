package helpers

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type RandomUserDataProvider struct {
	firstNameFile  string
	lastNameFile   string
	cityFile       string
	textFile       string
	firstNameSlice []string
	lastNameSlice  []string
	textSlice      []string
	citySlice      []string
	prefix         string
}

var genderValues = []string{
	"male",
	"female",
}

func CreateRandomUserDataProvider(prefix string) (*RandomUserDataProvider, error) {
	randomUserDataProvider := &RandomUserDataProvider{
		firstNameFile: "data/first_names.all.txt",
		lastNameFile:  "data/last_names.all.txt",
		cityFile:      "data/world-cities.csv",
		textFile:      "data/bbc-text.csv",
		prefix:        prefix,
	}

	csvFirstName, err := os.Open(randomUserDataProvider.firstNameFile)
	defer csvFirstName.Close()
	if err != nil {
		return nil, fmt.Errorf("couldn't open the csv file: %s", err.Error())
	}
	r := csv.NewReader(csvFirstName)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read first names: %s", err.Error())
		}
		randomUserDataProvider.firstNameSlice = append(randomUserDataProvider.firstNameSlice, strings.Title(strings.ToLower(record[0])))
	}

	csvLastName, err := os.Open(randomUserDataProvider.lastNameFile)
	defer csvLastName.Close()
	if err != nil {
		return nil, fmt.Errorf("couldn't open the csv file: %s", err.Error())
	}
	r = csv.NewReader(csvLastName)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read first names: %s", err.Error())
		}
		randomUserDataProvider.lastNameSlice = append(randomUserDataProvider.lastNameSlice, strings.Title(strings.ToLower(record[0])))
	}
	csvText, err := os.Open(randomUserDataProvider.textFile)
	defer csvText.Close()
	if err != nil {
		return nil, fmt.Errorf("couldn't open the csv file: %s", err.Error())
	}
	r = csv.NewReader(csvText)
	header := false
	for {
		record, err := r.Read()
		if !header {
			header = true
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read first names: %s", err.Error())
		}
		runes := []rune(record[1])
		randomUserDataProvider.textSlice = append(randomUserDataProvider.textSlice, string(runes[0:int(math.Min(512, float64(len(runes))))]))
	}
	csvCity, err := os.Open(randomUserDataProvider.lastNameFile)
	defer csvCity.Close()
	if err != nil {
		return nil, fmt.Errorf("couldn't open the csv file: %s", err.Error())
	}
	r = csv.NewReader(csvCity)
	header = false
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read first names: %s", err.Error())
		}
		randomUserDataProvider.citySlice = append(randomUserDataProvider.citySlice, strings.Title(strings.ToLower(record[0])))
	}

	return randomUserDataProvider, nil
}

func (dp *RandomUserDataProvider) FirstName() string {
	return dp.firstNameSlice[rand.Intn(len(dp.firstNameSlice))]
}

func (dp *RandomUserDataProvider) LastName() string {
	return dp.lastNameSlice[rand.Intn(len(dp.lastNameSlice))]
}

func (dp *RandomUserDataProvider) City() string {
	return dp.citySlice[rand.Intn(len(dp.citySlice))]
}

func (dp *RandomUserDataProvider) Text() string {
	return dp.textSlice[rand.Intn(len(dp.textSlice))]
}

func (dp *RandomUserDataProvider) Age() string {
	return strconv.Itoa(rand.Intn(108) + 1)
}

func (dp *RandomUserDataProvider) Gender() string {
	return genderValues[rand.Intn(2)]
}

func (dp *RandomUserDataProvider) Login(counter int) string {
	return dp.prefix + strconv.Itoa(counter)
}
