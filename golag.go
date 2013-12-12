package main

import (
  "github.com/codegangsta/martini"
  "html/template"
  "math/rand"
  "net/http"
  "time"
  "strconv"
)

type Post struct {
  Title     string
  Timestamp time.Time
  Content   template.HTML
}

type Site struct {
  Title string
  Posts []Post
}

var site Site

func main() {
  site = Site{
    Title: "Test Golag",
    Posts: make([]Post, 0),
  }
  makeFakePosts()

  m := martini.Classic()
  m.Get("/", index)

  m.Run()
}

func index(w http.ResponseWriter) {
  /*temp := Post{
    Title:     "Test",
    Timestamp: time.Now(),
    Content:   "<p>Hello <em>World</em>!</p>",
  }*/

  templ, _ := template.ParseFiles("templates/index.html")
  templ.Execute(w, site)
}

func makeFakePosts() {
  rand.Seed(time.Now().Unix())
  numPosts := rand.Intn(13) + 2

  for i := 1; i <= numPosts; i++ {
    p := Post{
      Title:     "Post " + strconv.Itoa(i),
      Timestamp: time.Now(),
      Content:   template.HTML("Content in <code>post #" + strconv.Itoa(i) + "</code>"),
    }

    site.Posts = append(site.Posts, p)
  }
}
