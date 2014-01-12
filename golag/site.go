package golag

import (
  "io/ioutil"
  "github.com/howeyc/fsnotify"
  "log"
  "sort"
  "html/template"
  "path/filepath"
)

type Site struct {
  Title string
  Posts []*Post
  Config map[string]string
  PostTemplate *template.Template
  IndexTemplate *template.Template
  Templates *template.Template
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

  sort.Sort(sort.Reverse(sortablePost(posts)))
  return posts, nil
}

func (s *Site) ReadTemplatesDirectory() (*template.Template, error) {
  fis, err := ioutil.ReadDir(s.Config["templateDir"])
  if err != nil {
    return nil, newError("Couldn't read templates directory", err)
  }
  
  files := make([]string, 0)
  for _, fi := range fis {
    files = append(files, filepath.Join(s.Config["templateDir"], fi.Name()))

    if err != nil {
      return nil, newError("Couldn't read template `" + fi.Name() + "`", err)
    }
  }

  templs, err := template.ParseFiles(files...)
  if err != nil {
    log.Println("\n\nError!", err, "\n\n")
    return nil, err
  }

  return template.Must(templs, nil), nil
}

func (s *Site) FindPost(slug string) (*Post) {
  for _, post := range s.Posts {
    if post.Slug == slug {
      return post
    }
  }

  return nil
}

func (s *Site) WatchPostChanges() {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Println("Couldn't make an fswatcher", err)
    return
  }

  done := make(chan bool)

  go func() {
    for {
      select {
      case <- watcher.Event:
        posts, err := s.ReadPostDirectory()
        if err != nil {
          log.Println("Couldn't refresh posts")
          continue
        }

        s.Posts = posts
      case <- watcher.Error:
        continue
      }
    }
  }()

  err = watcher.Watch(s.Config["postDir"])
  if err != nil {
    log.Println("Couldn't watch postDir", err)
    return
  }

  <- done
  watcher.Close()
}

type sortablePost []*Post
func (s sortablePost) Len() int { return len(s) }
func (s sortablePost) Less(i, j int) bool { return s[i].Timestamp.Before(s[j].Timestamp) }
func (s sortablePost) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
