package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getReviews(id int64) ([]Review, error) {
	var reviews []Review
	res, err := http.Get(fmt.Sprintf("http://localhost:4004/product/%d/review", id))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("here2")
		return nil, err
	}
	if err = json.Unmarshal(body, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func getUser(id int64) (string, error) {
	res, err := http.Get(fmt.Sprintf("http://localhost:4000/users/%d", id))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("here2")
		return "", err
	}
	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err = json.Unmarshal(body, &user); err != nil {
		return "", err
	}
	return user.Name, nil
}
