package main

import (
  "github.com/codegangsta/martini"
  "html/template"
  "net/http"
  "./golag"
  "fmt"
)

var site *golag.Site

func main() {
  site = &golag.Site{
    Title: "Test Golag",
    Posts: make([]*golag.Post, 0),
    Config: make(map[string]string),
  }

  site.Config["postDir"] = "content"

  posts, err := site.ReadPostDirectory() 
  if err != nil {
    fmt.Println("Couldn't read post directory")
    fmt.Println(err)
  }
  site.Posts = posts

  m := martini.Classic()
  m.Get("/", index)
  m.Get("/post/:slug", post)

  m.Run()
}

func index(w http.ResponseWriter) {
  templ, _ := template.ParseFiles("templates/index.html")
  templ.Execute(w, site)
}

func post(params martini.Params, w http.ResponseWriter) {
  post := site.FindPost(params["slug"])
  if post == nil {
    // TODO: make 404 better
    fmt.Fprintf(w, "No post found")
    return
  }

  templ, _ := template.ParseFiles("templates/post.html")
  templ.Execute(w, post)
}
