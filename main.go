package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

func loadPage(color string) (*Page, error) {
	filename := "test.txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	title := color
	return &Page{Title: title, Body: body}, nil
}

// Обработчик запросов. Принимает поток вывода от пакета "net/http",
// а также информацию о запроса
func handler(w http.ResponseWriter, r *http.Request) {
	println("inter from html")
	color := "green"
	font := r.FormValue("fontList")
	inputText := r.FormValue("inputText")
	println(" input text " + inputText)
	if r.FormValue("textcolor") != "" {
		color = r.FormValue("textcolor")
	}
	background := r.FormValue("bgcolor")
	if r.FormValue("bgcolor") != "" {
		background = r.FormValue("bgcolor")
	}
	//println(color)
	//println(background)
	templateFont := readFileFont(font + ".txt")
	f, err := os.Create("test.txt")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fcss, err := os.Create("./lib/main.css")
	err = fcss.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	writeFileFont(inputText, templateFont, "test.txt")
	writeFileCSS(color, background, "main.css")

	p, err := loadPage(color)
	if err != nil {
		println(" ошибка в загрузке файла")
		http.Redirect(w, r, "/err.html", http.StatusFound)
		return
	}

	var templates = r.URL.Path[1:]

	if templates == "" || templates == "/" {
		templates = "index.html"
	}

	var lGet = r.URL.Path[1:]      // Упростим напсание кода
	if lGet == "" || lGet == "/" { // Дабы разрешить "/" запросы
		lGet = "index.html"
	}

	//lData := readFile(lGet) // Считываем файл

	renderTemplate(w, r, p)

	// Отправляет файл в ответ на запрос
	//fmt.Fprintln(w, lData)
}

var templates = template.Must(template.ParseFiles("index.html"))

func renderTemplate(w http.ResponseWriter, r *http.Request, p *Page) {
	err := templates.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//templates.ExecuteTemplate(w, "err.html", p)
		http.Redirect(w, r, "/err.html", 404)
	}
}

func readFileFont(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var template []string
	for scanner.Scan() {
		template = append(template, scanner.Text())
	}
	return template
}

func writeFileCSS(c string, bc string, fcssname string) {

	fcss, err := os.OpenFile("./lib/"+fcssname, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fcss.WriteString("#art { " +
		"color: " + c + "; " +
		"background-color: " + bc + "; font-family: revert; " +
		"} ")

	if err != nil {
		fmt.Println(err)
		fcss.Close()
		return
	}
	err = fcss.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func writeFileFont(args string, template []string, fname string) {
	islineW := false
	num := 0

	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 1; i < 9; i++ {
		_, err = f.WriteString("\n")
		for j := range args {
			if args[j] == '\\' && args[j+1] == 'n' {
				islineW = true
				num = j
				break
			}

			for ascii := 32; ascii <= 126; ascii++ {
				if rune(args[j]) == rune(ascii) {
					_, err = f.WriteString(template[(ascii-32)*9+i])
				}

			}
		}
		fmt.Print()
	}

	if islineW {
		writeFileFont(args[num+2:len(args)], template, fname)
	}

	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readFile(iFileName string) string {
	// Считываем файл
	lData, err := ioutil.ReadFile(iFileName)
	var lOut string // Объявляем строчную переменную
	// Если файл существует - записываем его в переменную lOut
	if !os.IsNotExist(err) {
		lOut = string(lData)
	} else { // Иначе - отправляем стандартный текст
		lOut = "404"
	}
	return lOut // Возвращаем полученную информацию
}

// Наша главная функция
func main() {
	// При получени запроса к "/*", если не задано других обработчиков для данного
	// запроса, вызываем функцию "handler".
	http.HandleFunc("/", handler)
	http.HandleFunc("/err.html", handler)

	http.HandleFunc("/lib/main.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "lib/main.css")
	})

	http.HandleFunc("/lib/mainJS.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "lib/mainJS.js")
	})

	// Ну а как-же без этого?)
	log.Println("Запускаемся. Слушаем порт 8080")

	// Сканируем запросы к порту 8080. При наличии таковых - отвечаем так, как
	// указано выше
	http.ListenAndServe(":8080", nil)
}
