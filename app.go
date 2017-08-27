//TODO
// get img from it
// add https://tproger.ru/translations/backend-web-development/

package test_GoRESTful_API

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"runtime"
	"path"
	"fmt"
	"strings"
	"os"
)

var (
	dao     = memesDAO{}
	srv     = &http.Server{}
	appPath string
    mySigningKey = []byte("secret")
)

func init() {
	log.Print("Start init")
	defer log.Print("End init")

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	appPath = path.Dir(filename)
	fullPath := appPath + "/db.json"
	log.Print(fullPath, " config file")
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

func TestRoutes(router *mux.Router) {
	r := router.PathPrefix("/test").Subrouter()

	r.HandleFunc("/", testVar)
	r.HandleFunc("/products", testVar).Methods("POST")
	r.HandleFunc("/articles", testVar).Methods("GET")
	r.HandleFunc("/articles/{id}", testVar).Methods("GET", "PUT")
}

func RestRoutes(router *mux.Router) {
	api_router := router.PathPrefix("/api/").Subrouter()
	api_router.HandleFunc("/memes", getMemes).Methods("GET")
	api_router.HandleFunc("/memes", postMeme).Methods("POST")
	api_router.HandleFunc("/memes/{id}", getMemeID).Methods("POST")
}

func JWTRoutes(router *mux.Router) {
	router.HandleFunc("/get-token", GetToken).Methods("GET")
}

func Walker(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		p, err := route.GetPathRegexp()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}

		fmt.Println("Methodes:", strings.Join(m, ","))
		fmt.Println("Template", t)
		fmt.Println("RegExp", p)
		fmt.Println("")
		return nil
	})
}

func Start() {
	log.Print("Starting-------------------------------------------")
	router := mux.NewRouter()


	router.Handle("/", http.FileServer(http.Dir(appPath+"/views/")))

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(appPath+"./static/"))))
	router.Handle("/status", StatusHandler)
	router.Handle("/products", jwtMW.Handler(ProductsHandler))
	router.Handle("/products/{slug}/feedback",  jwtMW.Handler(AddFeedbackHandler))

	RestRoutes(router)
	TestRoutes(router)
	JWTRoutes(router)

	srv = &http.Server{Addr: ":3000", Handler: handlers.LoggingHandler(os.Stdout, router)}
	Walker(router)

	go func() {
		log.Print("Started")
		defer log.Print("Stoped")
		if err := srv.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()
}

func Stop() error {
	log.Print("Stopping...")
	return srv.Shutdown(nil)
}
