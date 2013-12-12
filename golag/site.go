package golag

import (
  "io/ioutil"
)

type Site struct {
  Title string
  Posts []*Post
  Config map[string]string
}

func (s *Site) ReadPostDirectory() ([]*Post, error) {
  posts := make([]*Post, 0)

  fis, err := ioutil.ReadDir(s.Config["postDir"])
  if err != nil {
    return nil, newError("Couldn't read posts directory", err)
  }

  for _, fi := range fis {
    post, err := ReadPost(fi, s)
    if err != nil {
      return nil, newError("Couldn't read post `" + fi.Name() + "`", err)
    }

    posts = append(posts, post)
  }

  return posts, nil
}

func (s *Site) FindPost(slug string) (*Post) {
  for _, post := range s.Posts {
    if post.Slug == slug {
      return post
    }
  }

  return nil
}