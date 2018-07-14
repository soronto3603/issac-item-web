package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Page struct {
	Title string
	Body  []byte
}

// const{

// }

func (p *Page) save() error {
	filename := p.Title + ".html"
	return ioutil.WriteFile(filename, p.Body, 0600)
}
func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// func main() {
// 	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
// 	// p1.save()
// 	p2, _ := loadPage("TestPage")
// 	fmt.Println(string(p2.Body))
// }

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("./static/templates/view.html")
	t.Execute(w, p)
}
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "./static/templates/edit", p)
}
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
func adminHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/admin/"):]
	p, _ := loadPage(title)
	t, _ := template.New("").ParseFiles("./static/templates/admin.html")
	t.Execute(w, p)
}
func issacInsertHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pf := r.PostForm

	// fmt.Println("[Server] ", pf["name"])
	fmt.Println("[Server] Receive data <=", r.PostForm)

	fmt.Printf("Server] mysql cheking start:)\n")
	db, err := sql.Open("mysql", "soronto:2262552a@/issac")
	if err != nil {
		fmt.Printf("Server] mysql load fail\n")
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO issac_item(no,id,name,description,img,battery,quote,effects,notes,synergies) VALUES(NULL,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()
	fmt.Printf("Server] Data Insert\n")

	// 기괴하고 해괴망측한 고랭
	// fmt.Println(reflect.TypeOf(pf["name"]))==> []string

	_, err = stmtIns.Exec(1, mapStringParse(pf, "name"), mapStringParse(pf, "description"), mapStringParse(pf, "img"), mapStringParse(pf, "battery"), mapStringParse(pf, "quote"), mapStringParse(pf, "effects"), mapStringParse(pf, "notes"), mapStringParse(pf, "synergies"))
	if err != nil {
		panic(err.Error())
	}
}

func mapStringParse(m url.Values, s string) string {
	fuckingString := strings.Join(m[s], "")
	return fuckingString
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/issac/insert/", issacInsertHandler)
	fmt.Printf("Server Start:8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
