package test_GoRESTful_API

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	. "github.com/Smoreg/test_GoRESTful_API/model"
	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
	"time"
	"log"
	"github.com/auth0/go-jwt-middleware"
)

func getMemes(w http.ResponseWriter, _ *http.Request) {
	memes, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, memes)
}
func getMemeID(w http.ResponseWriter, r *http.Request) {
	meme, err := dao.FindById(mux.Vars(r)["id"])

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, meme)

}
func postMeme(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var meme Meme
	if err := json.NewDecoder(r.Body).Decode(&meme); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	meme.ID = bson.NewObjectId()
	if err := dao.Insert(meme); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, meme)
}
func testVar(w http.ResponseWriter, r *http.Request) {

	a := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.Path + "\n"))
	tmp, _ := json.Marshal(a)
	w.Write(tmp)
}

func GetToken(w http.ResponseWriter, _ *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = "Vasia Pupkin"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Panic(err)
	}
	w.Write([]byte(tokenString))
}

var jwtMW = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,

})
