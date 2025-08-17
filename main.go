package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/OTakumi/chat/trace"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, r)
}

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading dotenv file")
	}
}

func main() {
	// dotenvファイルを読み込み
	Env_load()

	clientId := os.Getenv("CLIENT_ID")
	securityKey := os.Getenv("SECURITY_KEY")
	log.Println("Client Id: " + clientId)
	log.Println("Security Key: " + securityKey)

	addr := flag.String("addr", ":8080", "The address of the application.")
	flag.Parse()

	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey(securityKey)
	gomniauth.WithProviders(
		github.New("client id", "security key", "http://localhost:8080/auth/callback/github"),
		google.New(clientId, securityKey, "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// Start the room
	go r.run()

	// Start the web server
	log.Println("Starting the web server on this port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
