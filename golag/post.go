package golag

import (
  "time"
  "html/template"
  "os"
  "path/filepath"
  "bufio"
  "strings"
  "fmt"
  "regexp"
  "bytes"
  "github.com/russross/blackfriday"
  "strconv"
)

type Post struct {
  Title     string
  Timestamp time.Time
  Content   template.HTML
  Frontmatter map[string]string
  Slug string
}



func ReadPost(fileInfo os.FileInfo, s *Site) (*Post, error) {
  file, err := os.Open(filepath.Join(s.Config["postDir"], fileInfo.Name()))
  if err != nil {
    return nil, newError("Couldn't open file", err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  frontmatter, err := readFrontmatter(scanner)
  if err != nil {
    return nil, newError("Error reading frontmatter of `" + fileInfo.Name() + "`", err)
  }

  slug := getSlug(fileInfo.Name())
  //pubTime, err := time.Parse("2013-Jan-30", frontmatter["date"])
  //if err != nil {
  //  pubTime = fileInfo.ModTime()
  //}

  pubTime := fileInfo.ModTime()
  timeSplit := strings.Split(frontmatter["date"], "-")
  if len(timeSplit) >= 3 {
    year, errY := strconv.Atoi(timeSplit[0])
    month, errM := strconv.Atoi(timeSplit[1])
    day, errD := strconv.Atoi(timeSplit[2])

    loc, errL := time.LoadLocation("Local")

    if errY == nil && errM == nil && errD == nil && errL == nil {
      pubTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
    }
  }


  buffer := bytes.NewBuffer(nil)
  for scanner.Scan() {
    buffer.WriteString(scanner.Text() + "\n")
  }
  if err = scanner.Err(); err != nil {
    return nil, newError("Couldn't read file's contents `" + fileInfo.Name() + "`", err)
  }

  content := blackfriday.MarkdownCommon(buffer.Bytes())
  post := &Post{
    Title: frontmatter["title"],
    Timestamp: pubTime,
    Content: template.HTML(content),
    Frontmatter: frontmatter,
    Slug: slug,
  }

  return post, nil
}

func readFrontmatter(s *bufio.Scanner) (map[string]string, error) {
  m := make(map[string]string)

  inFrontmatter := false
  for s.Scan() {
    line := strings.Trim(s.Text(), " ")

    if line == "---" {
      if inFrontmatter {
        return m, nil
      } else {
        inFrontmatter = true
      }
    } else if inFrontmatter {
      segments := strings.SplitN(line, ":", 2)
      if len(segments) != 2 {
        return nil, fmt.Errorf("Invalid frontmatter")
      }
      m[segments[0]] = strings.Trim(segments[1], " ")
    } else if line == "" {
      return nil, fmt.Errorf("No frontmatter")
    }
  }

  if err := s.Err(); err != nil {
    return nil, newError("Scanner error.", err)
  }

  return nil, fmt.Errorf("Empty post...")
}

var slugRegex = regexp.MustCompile(`[^a-zA-Z\-0-9]`)
func getSlug(filename string) string {
  return slugRegex.ReplaceAllString(strings.Replace(filename, filepath.Ext(filename), "", 1), "-")
}