package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/garyburd/redigo/redis"
	"github.com/grayj/go-json-rest-middleware-tokenauth"
)

var (
	redisPool      *redis.Pool
	port           = os.Getenv("PORT")
	redisServer    = os.Getenv("REDIS_SERVER")
	redisPassword  = os.Getenv("REDIS_PASSWORD")
	tokenNamespace = "agent:"
	authRealm      = "klouds-agent"
)

// Serve a json api
func Serve() {

	redisPool = newPool(redisServer, redisPassword)

	api := rest.NewApi()

	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login"
		},
		IfTrue: &tokenauth.AuthTokenMiddleware{
			Realm: authRealm,
			Authenticator: func(token string) string {
				rd := redisPool.Get()
				defer rd.Close()
				user, _ := redis.String(rd.Do("GET", tokenNamespace+tokenauth.Hash(token)))
				return user
			},
		},
		IfFalse: &rest.AuthBasicMiddleware{
			Realm: authRealm,
			Authenticator: func(user string, password string) bool {
				if user == "user" && password == "password" {
					return true
				}
				return false
			},
		},
	})

	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/login", login),
		rest.Get("/containers", list),
		rest.Get("/container/:id", inspect),
		rest.Post("/image/build", build),
		rest.Post("/image/run", create),
		rest.Get("/container/start/:id", start),
		rest.Get("/container/stop/:id", stop),
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	log.Println("Port", port)
	log.Fatal(http.ListenAndServe(":"+port, api.MakeHandler()))
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func login(w rest.ResponseWriter, r *rest.Request) {
	token, err := tokenauth.New()
	if err != nil {
		rest.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rd := redisPool.Get()
	defer rd.Close()
	_, err = rd.Do("SET", tokenNamespace+tokenauth.Hash(token), r.Env["REMOTE_USER"].(string), "EX", 604800)

	if err != nil {
		log.Panicln("Internal Server Error", err)
		rest.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteJson(map[string]string{
		"access_token": token,
	})
}

func list(w rest.ResponseWriter, r *rest.Request) {
	docker := NewDockerClient()
	containers := List(docker)
	w.WriteJson(containers)
}

func inspect(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	docker := NewDockerClient()
	info := Inspect(docker, id)
	w.WriteJson(info)
}

func start(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	docker := NewDockerClient()
	Start(docker, id)
	w.WriteJson(id)
}

func stop(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	docker := NewDockerClient()
	Stop(docker, id)
	w.WriteJson(id)
}

func create(w rest.ResponseWriter, r *rest.Request) {

	name := r.Form["name"][0]
	image := r.Form["image"][0]
	docker := NewDockerClient()
	Create(docker, name, image)
	w.WriteJson("")
}

func build(w rest.ResponseWriter, r *rest.Request) {

	repoName := r.Form["repoName"][0]
	context := r.Form["context"][0]
	docker := NewDockerClient()
	Build(docker, repoName, context)
	w.WriteJson("containers")
}
