package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func GetUserId(address, username, password string) (string, *http.Cookie, error) {
	BaseURLV1 := address + "/login"
	body := &CCXLogin{
		Login:    username,
		Password: password,
	}
	jsonAuth := bytes.NewBuffer([]byte{})

	if err := json.NewEncoder(jsonAuth).Encode(body); err != nil {
		return "", nil, err
	}
	req, err := http.NewRequest("POST", BaseURLV1, jsonAuth)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != 200 {
		le := &LoginError{}
		if err := json.NewDecoder(res.Body).Decode(le); err != nil {
			return "", nil, err
		}
		return "", nil, fmt.Errorf("service returned non 200 status code: %s", le.Error)
	}
	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == "sid" {
			cookie = c
			break
		}
	}
	if cookie == nil {
		return "", nil, errors.New("session cookie not found")
	}
	rb := &LoginResponse{}
	if err := json.NewDecoder(res.Body).Decode(rb); err != nil {
		return "", nil, err
	}
	return rb.ID, cookie, nil
}
