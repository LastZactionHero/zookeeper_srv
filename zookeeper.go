package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "sort"
  "os"
  "strings"
  "strconv"
  "io"
  "time"
)

const photoDir string = "./photos"

func newPhotoFilename() string {
  now := time.Now()
  return fmt.Sprintf("%d.jpg", now.Unix())
}

func latestPhotoHandler(w http.ResponseWriter, r *http.Request) {
  latestFile := latestFile()
  path := photoDir + "/" + latestFile.Name()

  fmt.Println(path)
  http.ServeFile(w, r, path)
}

func postPhotoHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("postPhotoHandler")
  err := r.ParseMultipartForm(100000)

  if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }

  m := r.MultipartForm
  fileHeader := m.File["file"][0]
  file, err := fileHeader.Open()

  dst, _ := os.Create(photoDir + "/" + newPhotoFilename())
  io.Copy(dst, file);
}

type ByTime []os.FileInfo

func parseFileTimestamp(filename string) int64 {
  timestamp := strings.Replace(filename, ".jpg", "", 1)
  timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
  return timestampInt
}

func (s ByTime) Len() int {
  return len(s)
}

func (s ByTime) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}

func (s ByTime) Less(i, j int) bool {
  timestamp_i, timestamp_j := parseFileTimestamp(s[i].Name()), parseFileTimestamp(s[j].Name())
  return timestamp_i > timestamp_j
}

func latestFile() os.FileInfo {
  files, _ := ioutil.ReadDir(photoDir)
  sort.Sort(ByTime(files))
  return files[0]
}


func main() {
  fmt.Println("Hello Zookeeper")

  http.HandleFunc("/latest_photo", latestPhotoHandler)
  http.HandleFunc("/post_photox", postPhotoHandler)

  http.ListenAndServe(":8080", nil)
}