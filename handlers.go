package test_GoRESTful_API

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	. "github.com/Smoreg/test_GoRESTful_API/model"
	"fmt"
)


func getMemes(w http.ResponseWriter, r *http.Request){
	fmt.Println("\n--------------GET MEEEEMEEEES------------------\n")
	fmt.Println(r)
	memes, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, memes)

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