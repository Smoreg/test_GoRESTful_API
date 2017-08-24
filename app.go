package test_GoRESTful_API

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
	"runtime"
	"path"
)

var (
	dao = memesDAO{}
	srv = &http.Server{}
)

func init() {
	log.Print("Start init")
	defer log.Print("End init")

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	fullPath:=path.Dir(filename) + "/db.json"
	log.Print(fullPath, "config file")
	raw, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(raw, &dao)
	dao.connect()
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/memes", getMemes).Methods("GET")
	router.HandleFunc("/memes", postMeme).Methods("POST")
	//router.HandleFunc("/memes", PutMeme).Methods("PUT")
	//router.HandleFunc("/memes", DeleteMemes).Methods("DELETE")
	//router.HandleFunc("/memes/{id}", GetMemeID).Methods("GET")
	//router.HandleFunc("/memes/{id}", DeleteMemeID).Methods("DELETE")

	srv = &http.Server{Addr: ":3000", Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()
}

func Stop() error {
	return srv.Shutdown(nil)
}
