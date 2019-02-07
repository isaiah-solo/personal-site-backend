package main

import (
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "net/http"
  "reflect"
  "time"
)

type NullInt sql.NullInt64

type About struct {
  Jobs []Job `json:"jobs"`
}

type AboutResponse struct {
  About About `json:"about"`
}

type Icon struct {
  Name string `json:"name"`
  Website string `json:"website"`
}

type Job struct {
  Company string `json:"company"`
  Details []string `json:"details"`
  EndDate NullInt `json:"endDate"`
  Position string `json:"position"`
  Skills []Skill `json:"skills"`
  StartDate int `json:"startDate"`
  Website string `json:"website"`
}

type Profile struct {
  Headline string `json:"headline"`
  Icons []Icon `json:"icons"`
  Name string `json:"name"`
}

type ProfileResponse struct {
  Profile Profile `json:"profile"`
}

type Skill struct {
  Label string `json:"label"`
  Link string `json:"link"`
}

var username = "root"
var password = "admin"
var database = "site"

func (nullInt *NullInt) Scan(value interface{}) error {
  var i sql.NullInt64
  if err := i.Scan(value); err != nil {
    return err
  }
  if reflect.TypeOf(value) == nil {
    *nullInt = NullInt{i.Int64, false}
  } else {
    *nullInt = NullInt{i.Int64, true}
  }
  return nil
}

func (nullInt *NullInt) MarshalJSON() ([]byte, error) {
  if !nullInt.Valid {
    return []byte("null"), nil
  }
  return json.Marshal(nullInt.Int64)
}

func establishConnection() *sql.DB {
  db, err := sql.Open(
    "mysql",
    fmt.Sprintf(
      "%s:%s@/%s",
      username,
      password,
      database,
    ),
  )
  if err != nil {
    panic(err)
  }
  return db
}

func performQuery(db *sql.DB, query string, args ...interface{}) *sql.Rows {
  rows, err := db.Query(query, args...)
  if err != nil {
    panic(err)
  }
  return rows
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
  db := establishConnection()
  defer db.Close()
  rows := performQuery(
    db,
    "SELECT company, end_date, id, position, start_date, website FROM job",
  )
  about := About{
    Jobs: []Job{},
  }
  for rows.Next() {
    var job Job
    var jobID int
    if err := rows.Scan(
      &job.Company,
      &job.EndDate,
      &jobID,
      &job.Position,
      &job.StartDate,
      &job.Website,
    ); err != nil {
      panic(err)
    }
    rows := performQuery(
      db,
      "SELECT text FROM job_detail WHERE job_id = ?",
      jobID,
    )
    for rows.Next() {
      var detail string
      err := rows.Scan(
        &detail,
      )
      if err != nil {
        panic(err)
      }
      job.Details = append(job.Details, detail)
    }
    rows = performQuery(
      db,
      "SELECT skill_id FROM job_to_skill WHERE job_id = ?",
      jobID,
    )
    for rows.Next() {
      var skillID int
      if err := rows.Scan(
        &skillID,
      ); err != nil {
        panic(err)
      }
      rows := performQuery(
        db,
        "SELECT label, link FROM skill WHERE id = ?",
        skillID,
      )
      if rows.Next() != true {
        panic("Only 1 profile should exist.")
      }
      skill := Skill{}
      if err := rows.Scan(
        &skill.Label,
        &skill.Link,
      ); err != nil {
        panic(err)
      }
      job.Skills = append(job.Skills, skill)
    }
    about.Jobs = append(about.Jobs, job)
  }
  response, err := json.Marshal(
    AboutResponse{
      About: about,
    },
  )
  if err != nil {
    panic(err)
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write([]byte(response))
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
  db := establishConnection()
  defer db.Close()
  rows := performQuery(
    db,
    "SELECT headline, id, name FROM profile",
  )
  if rows.Next() != true {
    panic("Only 1 profile should exist.")
  }
  var profile Profile
  var profileID int
  if err := rows.Scan(
    &profile.Headline,
    &profileID,
    &profile.Name,
  ); err != nil {
    panic(err)
  }
  rows = performQuery(
    db,
    "SELECT name, website FROM profile_icon WHERE profile_id = ?",
    profileID,
  )
  for rows.Next() {
    var icon Icon
    if err := rows.Scan(
      &icon.Name,
      &icon.Website,
    ); err != nil {
      panic(err)
    }
    profile.Icons = append(profile.Icons, icon)
  }
  response, err := json.Marshal(
    ProfileResponse{
      Profile: profile,
    },
  )
  if err != nil {
    panic(err)
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write([]byte(response))
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/api/static/about", handleAbout).
    Methods("GET")
  r.HandleFunc("/api/static/profile", handleProfile).
    Methods("GET")
  srv := &http.Server{
    Addr:  "127.0.0.1:8081",
    Handler:  r,
    ReadTimeout:  15 * time.Second,
    WriteTimeout:  15 * time.Second,
  }
  if err := srv.ListenAndServe(); err != nil {
    panic(err)
  }
}
