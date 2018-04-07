package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	FileInfoURL string = "https://slack.com/api/files.info"
	UserInfoURL string = "https://slack.com/api/users.info"
)

type Helper struct {
	Token  string
	Client *http.Client
}

func DefaultHelper() *Helper {
	return &Helper{token, http.DefaultClient}
}

func (h *Helper) GetFileInfo(fileID string) (info map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", FileInfoURL, nil)
	if err != nil {
		return info, err
	}
	q := req.URL.Query()
	q.Add("token", h.Token)
	q.Add("file", fileID)
	req.URL.RawQuery = q.Encode()

	resp, err := h.Client.Do(req)
	if err != nil {
		return info, err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&info)
	if info["ok"].(bool) == false {
		err = errors.New(info["error"].(string))
		return info, err
	}

	return info["file"].(map[string]interface{}), err
}

func (h *Helper) GetUserInfo(userID string) (info map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", UserInfoURL, nil)
	if err != nil {
		return info, err
	}
	q := req.URL.Query()
	q.Add("token", h.Token)
	q.Add("user", userID)
	req.URL.RawQuery = q.Encode()

	resp, err := h.Client.Do(req)
	if err != nil {
		return info, err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&info)
	if info["ok"].(bool) == false {
		err = errors.New(info["error"].(string))
		return info, err
	}

	return info["user"].(map[string]interface{}), err
}

func (h *Helper) GetImageReader(imageURL string) (reader io.ReadCloser, err error) {
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+h.Token)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

func (h *Helper) Request(url string, data interface{}) (result map[string]interface{}, err error) {
	cleaned := make(map[string]interface{})
	cleaned["body"] = data
	cleaned["httpMethod"] = "POST"
	cleanedData, err := json.Marshal(cleaned)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(cleanedData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawData, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawData, &result)

	fmt.Println(result)
	return result, err
}
