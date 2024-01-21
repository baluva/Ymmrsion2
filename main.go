package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Aventurier struct {
	ID       int    `json:"id"`
	Nom      string `json:"nom"`
	Classe   string `json:"classe"`
	Niveau   int    `json:"niveau"`
	PointVie int    `json:"pointVie"`
}

var idCounter int = 0
var aventuriers []Aventurier

func main() {
	loadAventuriersFromJSON()
	css := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", css))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/profil", profilHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/modify", modifyHandler)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/home.html", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nom := r.FormValue("nom")
		classe := r.FormValue("classe")
		niveau := r.FormValue("niveau")
		pointVie := r.FormValue("pointVie")

		niveauInt := parseInt(niveau)
		pointVieInt := parseInt(pointVie)

		if nom != "" && classe != "" && niveauInt != -1 && pointVieInt != -1 {
			idCounter++
			aventurier := Aventurier{
				ID:       idCounter,
				Nom:      nom,
				Classe:   classe,
				Niveau:   niveauInt,
				PointVie: pointVieInt,
			}
			aventuriers = append(aventuriers, aventurier)
			saveAventuriersToJSON()
		}
	}

	renderTemplate(w, "templates/create.html", nil)
}

func profilHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/profil.html", aventuriers)
}

func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadAventuriersFromJSON() {
	file, err := ioutil.ReadFile("aventuriers.json")
	if err == nil {
		json.Unmarshal(file, &aventuriers)
	}
}

func saveAventuriersToJSON() {
	data, err := json.MarshalIndent(aventuriers, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile("aventuriers.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idParam := r.URL.Query().Get("id")
		idToDelete := parseInt(idParam)

		if idToDelete != -1 {
			index := findAventurierIndexByID(idToDelete)
			if index != -1 {
				aventuriers = append(aventuriers[:index], aventuriers[index+1:]...)
				saveAventuriersToJSON()

			}
		}
	}

	// Redirige vers la page de profil après la suppression
	http.Redirect(w, r, "/profil", http.StatusSeeOther)
}
func findAventurierIndexByID(id int) int {
	for i, aventurier := range aventuriers {
		if aventurier.ID == id {
			return i
		}
	}
	return -1
}
func modifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idParam := r.URL.Query().Get("id")
		idToModify := parseInt(idParam)

		var aventurierToModify *Aventurier
		for i, av := range aventuriers {
			if av.ID == idToModify {
				aventurierToModify = &aventuriers[i]
				break
			}
		}

		if aventurierToModify != nil {
			renderTemplate(w, "templates/modify.html", aventurierToModify)
		} else {
			http.Redirect(w, r, "/profil", http.StatusSeeOther)
		}
	} else if r.Method == http.MethodPost {
		idParam := r.FormValue("id")
		nom := r.FormValue("nom")
		classe := r.FormValue("classe")
		niveau := r.FormValue("niveau")
		pointVie := r.FormValue("pointVie")

		idToModify := parseInt(idParam)
		niveauInt := parseInt(niveau)
		pointVieInt := parseInt(pointVie)

		if idToModify != -1 {
			index := findAventurierIndexByID(idToModify)
			if index != -1 {
				aventuriers[index].Nom = nom
				aventuriers[index].Classe = classe
				aventuriers[index].Niveau = niveauInt
				aventuriers[index].PointVie = pointVieInt

				saveAventuriersToJSON()
			}
		}

		http.Redirect(w, r, "/profil", http.StatusSeeOther)
	}
}
