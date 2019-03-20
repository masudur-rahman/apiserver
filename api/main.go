package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"gopkg.in/macaron.v1"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var engine *xorm.Engine

type Worker struct {
	Username string `json:"username" xorm:"pk not null unique"`

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`

	City     string `json:"city"`
	Division string `json:"division"`

	Position string `json:"position"`
	Salary   int64  `json:"salary"`

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	Version   int       `xorm:"version"`
}

// List of workers and authenticated users
var Workers []Worker
var authUser = make(map[string]string)

var srvr http.Server
var byPass bool = true
var stopTime int16

func StartXormEngine() {
	var err error
	connStr := "user=masud password=masud123 host=127.0.0.1 port=5432 dbname=apiserver sslmode=disable"

	engine, err = xorm.NewEngine("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	logFile, err := os.Create("apiserver.log")
	if err != nil {
		log.Println(err)
	}
	logger := xorm.NewSimpleLogger(logFile)
	logger.ShowSQL(true)
	engine.SetLogger(logger)

	if engine.TZLocation, err = time.LoadLocation("Asia/Dhaka"); err != nil {
		log.Println(err)
	}
}

// Handler Functions....

func Welcome(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode("Congratulations...! Your API Server is up and running... :) "); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func WelcomeToAppsCode(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode("Welcome to AppsCode Ltd.. Available Links are : `/appscode/workers`, `/appscode/workers/{username}`"); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ShowAllWorkers(w http.ResponseWriter, r *http.Request) {

	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte(info)); err != nil {
			log.Println(err)
		}
		return
	}

	var workers []Worker
	if err := engine.Find(&workers); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(workers); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func ShowSingleWorker(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte(info)); err != nil {
			log.Println(err)
		}
		return
	}

	worker := new(Worker)
	worker.Username = ctx.Params("username")
	exist, err := engine.Get(worker)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if !exist {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("404 - Content Not Found")); err != nil {
			log.Println(err)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(worker); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func AddNewWorker(w http.ResponseWriter, r *http.Request) {

	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte(info)); err != nil {
			log.Println(err)
		}
		return
	}

	var worker Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error decoding provided data")); err != nil {
			log.Println(err)
		}
		return
	}

	if worker.Username == "" {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Username must be provided")); err != nil {
			log.Println(err)
		}
		return
	}

	newWorker := new(Worker)
	newWorker.Username = worker.Username
	if exist, _ := engine.Get(newWorker); exist {
		w.WriteHeader(http.StatusConflict)
		if _, err := w.Write([]byte("409 - username already exists")); err != nil {
			log.Println(err)
		}
		return
	}

	// Check if it exists in deleted accounts
	newWorker = new(Worker)
	newWorker.Username = worker.Username
	if exist, _ := engine.Unscoped().Get(newWorker); exist {
		w.WriteHeader(http.StatusConflict)
		if _, err := w.Write([]byte("409 - username already exists")); err != nil {
			log.Println(err)
		}
		return
	}

	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := session.Insert(&worker); err != nil {
		if err = session.Rollback(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := session.Commit(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(worker); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateWorkerProfile(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {

	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte(info)); err != nil {
			log.Println(err)
		}
		return
	}

	worker := new(Worker)
	worker.Username = ctx.Params("username")
	exist, err := engine.Get(worker)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !exist {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("404 - Content Not Found")); err != nil {
			log.Println(err)
		}
		return
	}

	newWorker := new(Worker)
	if err := json.NewDecoder(r.Body).Decode(newWorker); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error decoding provided data")); err != nil {
			log.Println(err)
		}
		return
	}
	if newWorker.Username != worker.Username {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("405 - Username can't be changed")); err != nil {
			log.Println(err)
		}
		return
	}

	// Updated information assignment
	worker.FirstName = newWorker.FirstName
	worker.LastName = newWorker.LastName
	worker.City = newWorker.City
	worker.Division = newWorker.Division
	worker.Salary = newWorker.Salary

	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := session.ID(worker.Username).Update(worker); err != nil {
		log.Println(err)
		if err = session.Rollback(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := session.Commit(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("201 - Updated successfully")); err != nil {
		log.Println(err)
	}
	/*if err := json.NewEncoder(w).Encode(worker); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}*/
}

func DeleteWorker(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte(info)); err != nil {
			log.Println(err)
		}
		return
	}

	worker := new(Worker)
	worker.Username = ctx.Params("username")
	exist, err := engine.Get(worker)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !exist {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("404 - Content Not Found")); err != nil {
			log.Println(err)
		}
		return
	}

	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := session.ID(worker.Username).Delete(worker); err != nil {
		log.Println(err)
		if err = session.Rollback(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := session.Commit(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("200 - Deleted Successfully")); err != nil {
		log.Println(err)
	}
}

// Creating initial worker profiles
func CreateInitialWorkerProfile() {
	Workers = make([]Worker, 0)
	worker := Worker{
		Username:  "masud",
		FirstName: "Masudur",
		LastName:  "Rahman",
		City:      "Madaripur",
		Division:  "Dhaka",
		Position:  "Software Engineer",
		Salary:    55,
	}
	Workers = append(Workers, worker)

	worker = Worker{
		Username:  "fahim",
		FirstName: "Fahim",
		LastName:  "Abrar",
		City:      "Chittagong",
		Division:  "Chittagong",
		Position:  "Software Engineer",
		Salary:    55,
	}
	Workers = append(Workers, worker)

	worker = Worker{
		Username:  "tahsin",
		FirstName: "Tahsin",
		LastName:  "Rahman",
		City:      "Chittagong",
		Division:  "Chittagong",
		Position:  "Software Engineer",
		Salary:    55,
	}
	Workers = append(Workers, worker)

	worker = Worker{
		Username:  "jenny",
		FirstName: "Jannatul",
		LastName:  "Ferdows",
		City:      "Chittagong",
		Division:  "Chittagong",
		Position:  "Software Engineer",
		Salary:    55,
	}
	Workers = append(Workers, worker)

	if exist, _ := engine.IsTableExist(new(Worker)); !exist {
		if err := engine.CreateTables(new(Worker)); err != nil {
			log.Fatalln(err)
		}
	}

	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		log.Fatalln(err)
	}

	for _, user := range Workers {
		if _, err := session.Insert(&user); err != nil {
			if err = session.Rollback(); err != nil {
				log.Fatalln(err)
			}
		}
	}
	if err := session.Commit(); err != nil {
		log.Fatalln(err)
	}

	authUser["masud"] = "pass"
	authUser["admin"] = "admin"

}

func basicAuth(r *http.Request) (string, bool) {
	if byPass {
		return "", true
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "Error: Authorization Needed...!", false
	}

	authInfo := strings.SplitN(authHeader, " ", 2)

	userInfo, err := base64.StdEncoding.DecodeString(authInfo[1])

	if err != nil {
		return "Error: Error while decoding...!", false
	}
	userPass := strings.SplitN(string(userInfo), ":", 2)

	if len(userPass) != 2 {
		return "Error: Authorization failed...!", false
	}

	if pass, exist := authUser[userPass[0]]; exist {
		if pass != userPass[1] {
			return "Error: Unauthorized User", false
		} else {
			return "Success: Authorization Successful...!!", true
		}
	} else {
		return "Error: Unauthorized User...!", false
	}
}

func AssignValues(port string, bypass bool, stop int16) {
	srvr.Addr = ":" + port
	byPass = bypass
	stopTime = stop
}

func StartTheApp() {
	m := macaron.Classic()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	srvr.WriteTimeout = time.Second * 15
	srvr.ReadTimeout = time.Second * 15
	srvr.IdleTimeout = time.Second * 60

	srvr.Addr = "0.0.0.0:8080"

	srvr.Handler = m

	StartXormEngine()
	CreateInitialWorkerProfile()

	m.Get("/", Welcome)
	m.Group("appscode", func() {
		m.Get("/", WelcomeToAppsCode)
		m.Group("/workers", func() {
			m.Get("/", ShowAllWorkers)
			m.Get("/:username", ShowSingleWorker)
			m.Post("/", AddNewWorker)
			m.Put("/:username", UpdateWorkerProfile)
			m.Delete("/:username", DeleteWorker)
		})
	})

	log.Println("Starting the server")

	go func() {
		if err := srvr.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	// Channel to interrupt the server from keyboard
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, os.Interrupt)
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	//  Shutting down the server
	log.Println("Shutting down the server...!")

	time.Sleep(time.Second * time.Duration(stopTime))

	if err := srvr.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}

	log.Println("The server has been shut down...!")
	if err := engine.Close(); err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}
