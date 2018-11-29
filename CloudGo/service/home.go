package service

import (
	"net/http"
	"text/template"
	"time"

	"github.com/unrolled/render"
)

func apitest(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.Text(w, http.StatusOK,"The date now is: " +  time.Now().Format("2006/01/02 15:04:05"))
	}
}
func homeHandle(w http.ResponseWriter, r *http.Request) {
	//使用template.ParseFiles()实现模板的渲染输出
	//文件路径的根目录以可执行文件为基准
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}
