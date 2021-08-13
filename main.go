package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func home(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println("form:", req.Form)
	fmt.Println("path:", req.URL.Path)
	fmt.Println("scheme:", req.URL.Scheme)
	fmt.Println("long url:", req.Form["url_long"])
	for k, v := range req.Form {
		fmt.Println("key:", k, ",value:", strings.Join(v, ""))
	}
	fmt.Fprintln(w, "Hello This is go web")
}

func login(w http.ResponseWriter, req *http.Request) {
	fmt.Println("login method:", req.Method)
	if req.Method == "GET" {
		timestamp := strconv.Itoa(time.Now().Nanosecond())
		timeMD5 := md5.New()
		timeMD5.Write([]byte(timestamp))
		loginToken := fmt.Sprintf("%x", timeMD5.Sum(nil))
		t, _ := template.ParseFiles("login.gtpl")
		log.Println(t.Execute(w, loginToken))
	} else {
		req.ParseForm()
		loginToken := req.Form.Get("token")
		if loginToken == "" {
			log.Fatalln("login token error")
			return
		}

		fmt.Println("username:", template.HTMLEscapeString(req.Form.Get("username")))
		fmt.Println("password:", template.HTMLEscapeString(req.Form.Get("password")))
		template.HTMLEscape(w, []byte(req.Form.Get("username")))
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		curTime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curTime, 10))
		uploadToken := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, uploadToken)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handle, err := r.FormFile("uploadfile")
		if err != nil {
			log.Fatalln("from file error", err)
			return
		}

		defer file.Close()
		fmt.Fprintf(w, "%v", handle.Header)

		f, err := os.OpenFile("./file/"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln("open file error.", err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main_1() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":19191", nil)
	if err != nil {
		log.Fatalln("listen 19191 error")
	}
}

type News struct {
	Title   string `Title`
	Content string `Content`
	Author  string `Author`
}

func routeVue(w http.ResponseWriter, r *http.Request) {
	news := News{Title: "Learn Go Vue", Content: "Hello go Vue", Author: "Cooper"}
	t, err := template.ParseFiles("webapp/index.html")
	if err != nil {
		log.Fatalln("parse file vue.html failed!", err)
	}
	t.ExecuteTemplate(w, "index.html", news)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/vue", routeVue)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("webapp/static"))))
	err := http.ListenAndServe(":19191", nil)
	if err != nil {
		log.Fatalln("listen 19191 error")
	}
}
