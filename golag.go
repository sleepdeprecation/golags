package main

import (
  "github.com/codegangsta/martini"
  //"html/template"
  "net/http"
  "./golag"
  "fmt"
  "log"
)

var site *golag.Site

func main() {
  site = &golag.Site{
    Title: "Words, words, words",
    Posts: make([]*golag.Post, 0),
    Config: make(map[string]string),
  }

  site.Config["postDir"] = "content"
  site.Config["templateDir"] = "templates"
  site.Config["root"] = "/"

  posts, err := site.ReadPostDirectory() 
  if err != nil {
    fmt.Println("Couldn't read post directory")
    fmt.Println(err)
  }
  site.Posts = posts
  templs, err := site.ReadTemplatesDirectory()
  if err != nil {
    log.Fatal(err)
  }
  site.Templates = templs

  go site.WatchPostChanges()

  m := martini.Classic()
  m.Get("/", index)
  m.Get("/post/:slug", post)
  m.Run()
}

func index(w http.ResponseWriter) {
  err := site.Templates.ExecuteTemplate(w, "index.html", golag.GetPage(site, nil))
  if err != nil {
    log.Fatal("\n\nError!\n", err)
  }
}

func post(params martini.Params, w http.ResponseWriter) {
  post := site.FindPost(params["slug"])
  if post == nil {
    // TODO: make 404 better
    fmt.Fprintf(w, "No post found")
    return
  }

  err := site.Templates.ExecuteTemplate(w, "default", golag.GetPage(site, post))
  if err != nil {
    log.Println("\n\n\nError!", err)
  }
}
