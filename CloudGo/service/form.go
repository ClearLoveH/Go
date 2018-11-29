package service

import (
	"log"
	"net/http"
	"text/template"
)

// Login .
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["username"][0]
	pass := r.Form["password"][0]
	//请求的是登录，判断是否正确输入用户名及密码
	if len(name) == 0 || len(pass) == 0 {
		log.Fatal("Need username and password")
		http.Error(w, "Need username and password", 502)
	} else{
		t := template.Must(template.ParseFiles("templates/chart.html"))
		t.Execute(w, map[string]string{
			"Name": name,
			"Pass": pass,
		})
	}
}
