package main

import (
	"discord/internal/app"
	"discord/internal/handlers/user"
	"discord/internal/router"
	"discord/util"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)



func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
  coreApp := app.Init()
  userHandler := user.NewUserHandler(coreApp)

	r := router.SetupRouter(userHandler)
	port := util.GetEnv("PORT", "8080")
	log.Printf("Server listening on port %s...", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

// func (a *app.App)GetHTML(w http.ResponseWriter, r *http.Request) {
// 	title := chi.URLParam(r, "title")
// 	_, err := a.Q.CreateFoo(r.Context(),title)
//   if err != nil {
// 	slog.Warn("db err", "error", err)
//   }
// 	tmpl := template.Must(template.ParseFiles("index.html"))
// 	data := pageData{
// 		Title: fmt.Sprintf("<h1>%s</h1>",title),
// 		Items: []string{"xdd", "item2", "cs"},
// 	}
// 	tmpl.Execute(w, data)
// 	// w.Write([]byte("<h1>Hello WEB</h1>"))
// }
// func webRouter(r chi.Router) {
// 	r.Get("/title/{title}", app.GetHTML)
// }
