package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type LoginError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
type LoginResponse struct {
	ID               string `json:"id"`
	Login            string `json:"login"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Status           int    `json:"status"`
	EmailConfirmed   bool   `json:"emailConfirmed"`
	PendingRemoval   bool   `json:"pendingRemoval"`
	CompanyID        string `json:"companyId"`
	GroupID          string `json:"groupId"`
	AllowNewsletters bool   `json:"allowNewsletters"`
	TermsAccepted    struct {
		TermsAndConditionsV1 bool `json:"termsAndConditionsV1"`
		PrivacyPolicyV1      bool `json:"privacyPolicyV1"`
	} `json:"termsAccepted"`
	ExternalIds  interface{} `json:"externalIds"`
	Origin       string      `json:"origin"`
	DisableTrial bool        `json:"disableTrial"`
	Tasks        []int       `json:"tasks"`
}

func GetUserId(address string, username string, password string) (id string, seessionId *http.Cookie) {
	BaseURLV1 := address + "/login"
	body := &CCXLogin{
		Login:    username,
		Password: password,
	}
	jsonAuth := new(bytes.Buffer)
	json.NewEncoder(jsonAuth).Encode(body)
	req, _ := http.NewRequest("POST", BaseURLV1, jsonAuth)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	authClient := &http.Client{}
	res, err := authClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal(res.StatusCode)
	}
	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	cookie := res.Cookies()[0]
	if res.StatusCode == 500 {
		var LoginErrorResponse LoginError
		json.Unmarshal(responseBody, &LoginErrorResponse)
		log.Fatal(LoginErrorResponse.Error)
	}

	var CCXAuthResponse LoginResponse
	json.Unmarshal(responseBody, &CCXAuthResponse)
	return CCXAuthResponse.ID, cookie
}
