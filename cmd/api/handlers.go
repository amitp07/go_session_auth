package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"session-auth/internal/data"
	"session-auth/internal/dto"
	"session-auth/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
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

	err = app.models.CreateUserWithGroup(user, "")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := dto.UserResponse{
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

	otp := utils.GenerateOtp()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionId, err := app.config.redisClient.SetOtp(user.Username, otp)

	if err != nil {
		fmt.Printf("Error:mfa_redis_set: %s", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "mfa_session_id",
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		SameSite: http.SameSiteStrictMode,
	})

	w.Write([]byte("otp is: " + otp))

	// res, err := json.Marshal(user) if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// sessionId, err := app.config.redisClient.Set(user.Username)

	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "session_id",
	// 	Value:    sessionId,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(10 * time.Minute),
	// 	SameSite: http.SameSiteStrictMode,
	// })

	// w.Write(res)

}

func (app *application) verifyOtp(w http.ResponseWriter, r *http.Request) {
	otp := chi.URLParam(r, "otp")

	if otp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("mfa_session_id")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("cookie %s", cookie.Value)

	sessionOtp := app.config.redisClient.Get(cookie.Value)

	if sessionOtp != otp {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := app.config.redisClient.SetSession("user_name_one")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Expires:  time.Now().Add(60 * time.Minute),
		SameSite: http.SameSiteStrictMode,
		Value:    session,
		Path:     "/",
	})

	w.Write([]byte(fmt.Sprintf("session:%s", session)))
}
