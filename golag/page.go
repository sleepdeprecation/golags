package golag

import (
  "html/template"
  //"fmt"
)

type Page struct {
  Title string
  Type string
  Content template.HTML
  Post *Post
  Site *Site
  Timestamp timestamp
}

func GetPage(s *Site, p *Post) *Page {
  page := new(Page)
  page.Site = s

  if p == nil {
    page.Title = s.Title
    page.Type = "index"
  } else {
    page.Post = p
    page.Title = p.Title
    page.Content = p.Content
    page.Type = "post"
    page.Timestamp = p.Timestamp
  }

  //fmt.Println(page.Content)
  return page
}