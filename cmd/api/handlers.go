package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"session-auth/internal/data"
	"session-auth/internal/dto"
	"session-auth/internal/services"
	"session-auth/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
)

// handler to register new user
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

// Dummy handle to simulate protected route, this should be removed before moving to prod
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

// handler to take the sign-in requests
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

	redisVal, err := json.Marshal(dto.MfaSession{
		Username: user.Username,
		Otp:      otp,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionId, err := app.config.redisClient.SetOtp(redisVal)

	if err != nil {
		fmt.Printf("Error:mfa_redis_set: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "mfa_session_id",
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		SameSite: http.SameSiteStrictMode,
	})

	emailCfg := services.OtpEmailConfig()

	if err = emailCfg.Send(user.Email, otp); err != nil {
		fmt.Printf("could not send otp %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("otp send successfully."))

}

// handler to verify the OTP
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

	sessionOtpStr := app.config.redisClient.Get(cookie.Value)

	var sessionOtp dto.MfaSession
	err = json.Unmarshal([]byte(sessionOtpStr), &sessionOtp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if otp != sessionOtp.Otp {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = app.config.redisClient.Delete(cookie.Value)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session, err := app.config.redisClient.SetSession(sessionOtp.Username)

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
