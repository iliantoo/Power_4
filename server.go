package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Informations struct {
	Joueur1      string
	Joueur2      string
	Difficulte   string
	Grille       [][]int
	JoueurActuel int
	Gagnant      int
	NbCols       int
}

var currentGame *Informations

func createGrid(difficulty string) [][]int {
	var rows, cols int
	switch difficulty {
	case "6x7":
		rows, cols = 6, 7
	case "6x9":
		rows, cols = 6, 9
	case "7x8":
		rows, cols = 7, 8
	default:
		rows, cols = 6, 7
	}

	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
	}
	return grid
}

func dropToken(grid [][]int, col, joueur int) bool {
	for i := len(grid) - 1; i >= 0; i-- {
		if grid[i][col] == 0 {
			grid[i][col] = joueur
			return true
		}
	}
	return false
}

func checkWinner(grid [][]int, joueur int) bool {
	rows := len(grid)
	cols := len(grid[0])

	// Horizontal
	for r := 0; r < rows; r++ {
		for c := 0; c < cols-3; c++ {
			if grid[r][c] == joueur && grid[r][c+1] == joueur &&
				grid[r][c+2] == joueur && grid[r][c+3] == joueur {
				return true
			}
		}
	}

	// Vertical
	for c := 0; c < cols; c++ {
		for r := 0; r < rows-3; r++ {
			if grid[r][c] == joueur && grid[r+1][c] == joueur &&
				grid[r+2][c] == joueur && grid[r+3][c] == joueur {
				return true
			}
		}
	}

	// Diagonale /
	for r := 3; r < rows; r++ {
		for c := 0; c < cols-3; c++ {
			if grid[r][c] == joueur && grid[r-1][c+1] == joueur &&
				grid[r-2][c+2] == joueur && grid[r-3][c+3] == joueur {
				return true
			}
		}
	}

	// Diagonale \
	for r := 0; r < rows-3; r++ {
		for c := 0; c < cols-3; c++ {
			if grid[r][c] == joueur && grid[r+1][c+1] == joueur &&
				grid[r+2][c+2] == joueur && grid[r+3][c+3] == joueur {
				return true
			}
		}
	}

	return false
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./pages/index.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Println("Erreur template Index:", err)
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func Infos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if currentGame == nil || r.FormValue("col") == "" {
			j1 := r.FormValue("Joueur1")
			j2 := r.FormValue("Joueur2")
			diff := r.FormValue("difficulty")

			grid := createGrid(diff)

			currentGame = &Informations{
				Joueur1:      j1,
				Joueur2:      j2,
				Difficulte:   diff,
				Grille:       grid,
				JoueurActuel: 1,
				NbCols:       len(grid[0]),
			}
		} else {
			colStr := r.FormValue("col")
			col, _ := strconv.Atoi(colStr)
			dropToken(currentGame.Grille, col, currentGame.JoueurActuel)

			if checkWinner(currentGame.Grille, currentGame.JoueurActuel) {
				currentGame.Gagnant = currentGame.JoueurActuel
			} else {
				if currentGame.JoueurActuel == 1 {
					currentGame.JoueurActuel = 2
				} else {
					currentGame.JoueurActuel = 1
				}
			}
		}
	}

	tmpl, err := template.ParseFiles("./pages/info.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Println("Erreur template Info:", err)
		log.Fatal(err)
	}
	tmpl.Execute(w, currentGame)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/info", Infos)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
