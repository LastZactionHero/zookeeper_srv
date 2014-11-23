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
  "encoding/json"
)

const photoDir string = "./photos"
const capturingFilename string = "capture_status.txt"

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
  err := r.ParseMultipartForm(100000)

  if err != nil {
      fmt.Println(err.Error())
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }

  m := r.MultipartForm
  fileHeader := m.File["file"][0]
  file, err := fileHeader.Open()

  dst, _ := os.Create(photoDir + "/" + newPhotoFilename())
  io.Copy(dst, file);
}

func getCaptureHandler(w http.ResponseWriter, r *http.Request) {
  mapCapture := map[string]bool{"capturing": readCaptureStatus()}
  jsonCapture, _ := json.Marshal(mapCapture)

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, string(jsonCapture))
}

func postCaptureHandler(w http.ResponseWriter, r *http.Request) {
  capturing, _ := strconv.ParseBool(r.FormValue("capture"))
  storeCaptureStatus(capturing)

  mapCapture := map[string]bool{"capturing": capturing}
  jsonCapture, _ := json.Marshal(mapCapture)

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, string(jsonCapture))  
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

func storeCaptureStatus(capturing bool) bool { 
  c :=  []byte(strconv.FormatBool(capturing))
  err := ioutil.WriteFile(capturingFilename, c, 0644)
  return err != nil
}

func readCaptureStatus() bool {
  dat, err := ioutil.ReadFile(capturingFilename)
  if err != nil {
    return false
  }
  c, _ := strconv.ParseBool(string(dat))
  return c
}


func main() {  
  http.HandleFunc("/latest_photo", latestPhotoHandler)
  http.HandleFunc("/post_photo", postPhotoHandler)
  http.HandleFunc("/update_capture", postCaptureHandler)
  http.HandleFunc("/capture_status", getCaptureHandler)

  http.ListenAndServe(":8080", nil)
}