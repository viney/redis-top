package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var (
	Google = []byte("google")
	Baidu  = []byte("baidu")
	Bing   = []byte("bing")
)

var indexTemp = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := indexTemp.Execute(w, nil); err != nil {
		fmt.Println("indexHandler.indexTemp.Execute: ", err.Error())
		return
	}
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("topHandler.ParseForm: ", err.Error())
		return
	}

	member := r.FormValue("top")
	switch member {
	case "google":
		{
			f, err := Zincrby("top", Google, 1)
			if err != nil {
				fmt.Println("topHandler.Zincrby.Google: ", err.Error())
				return
			} else if f <= 0 {
				fmt.Println("topHandler.Zincrby.Google: ", f)
				return
			}
		}
	case "baidu":
		{
			f, err := Zincrby("top", Baidu, 1)
			if err != nil {
				fmt.Println("topHandler.Zincrby.Baidu: ", err.Error())
				return
			} else if f <= 0 {
				fmt.Println("topHandler.Zincrby.Baidu: ", f)
				return
			}
		}
	case "bing":
		{
			f, err := Zincrby("top", Bing, 1)
			if err != nil {
				fmt.Println("topHandler.Zincrby.Bing: ", err.Error())
				return
			} else if f <= 0 {
				fmt.Println("topHandler.Zincrby.Bing: ", f)
				return
			}
		}
	default:
		{
			panic("never not to here")
		}
	}

	// 获取所有成员
	bytesMap, err := Zrevrange("top", 0, -1)
	if err != nil {
		fmt.Println("topHandler.Zrang: ", err.Error())
	}

	var Tops = struct {
		Google int
		Baidu  int
		Bing   int
	}{}

	for member, score := range bytesMap {
		switch member {
		case string(Google):
			Tops.Google = int(score)
		case string(Baidu):
			Tops.Baidu = int(score)
		case string(Bing):
			Tops.Bing = int(score)
		default:
			panic("never not to here")
		}
	}

	if err := indexTemp.Execute(w, Tops); err != nil {
		fmt.Println("topHandler.indexTemp.Execute: ", err.Error())
		return
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/top", topHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
