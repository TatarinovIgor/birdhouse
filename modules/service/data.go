package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (service ATWalletService) RequestToDatabase() (io.ReadCloser, error) {
	client := http.Client{}
	requestType := "GET"
	url := "https://no-code-wallet.bird-house.org/version-test/api/1.1/obj/User?api_token=0c7c5967c58feaa9e93ee34e03e4cbc7"
	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		return nil, fmt.Errorf("can't make %s request for url: %s, err %v", requestType, url, err)
	}

	// make a request
	result, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't request %s for url: %s, err %v", requestType, url, err)
	}
	if result.StatusCode/100 != 2 {
		defer func() { _ = result.Body.Close() }()
		bodyBytes, _ := io.ReadAll(result.Body)
		return nil, fmt.Errorf("unexpected code from %s request for url: %s, code %v, message: %s",
			requestType, url, result.StatusCode, string(bodyBytes))
	}
	return result.Body, nil
}

func (service ATWalletService) GetUsersFromDatabase() (*[]UsersData, error) {
	result, err := service.RequestToDatabase()
	if err != nil {
		return nil, err
	}
	Users := BubbleUsersData{}
	err = json.NewDecoder(result).Decode(&Users)
	if err != nil {
		return nil, err
	}
	return &Users.Response.Result, nil
}
