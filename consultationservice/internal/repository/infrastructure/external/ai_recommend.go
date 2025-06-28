package external

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func RecommendationApi(dto.FilterRawData) (interface{}, error) {
	url := "https://api.chessy.dev/ai-recommend/therapist"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Print(err.Error())
		return nil, errors.New(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		return nil, errors.New(err.Error())
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print(err.Error())
		return nil, errors.New(err.Error())
	}
	var user model.ClientSwipes
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.ClientId)
	return body, nil
}
