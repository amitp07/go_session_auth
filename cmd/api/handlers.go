package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"session-auth/internal/data"
	"session-auth/internal/dto"
	"session-auth/internal/utils"
	"time"
)

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	var user dto.UserRequest
	encoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := encoder.Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := app.models.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := dto.UserResponse{
		ID:       id,
		Username: user.Username,
	}
	json, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(json)
}

func (app *application) getAllUsers(w http.ResponseWriter, r *http.Request) {

	var users []data.User
	app.models.GetAllUsers(&users)

	json, err := json.Marshal(users)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var userReq dto.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := app.models.GetUserByUsername(userReq.Username)

	if err != nil {
		fmt.Printf("Error %s\n", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := utils.MatchPassword(user.Password, userReq.Password); err != nil {
		fmt.Printf("password does not match %s\n", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionId, err := app.config.redisClient.Set(user.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		SameSite: http.SameSiteStrictMode,
	})

	w.Write(res)

}
