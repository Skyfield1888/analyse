package backend

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var RedirectTO string = "/table/test"

var Template = template.Must(
	template.New("table.html").Funcs(template.FuncMap{
		"untilCount": func(count int) []int {
			result := make([]int, count)
			for i := range result {
				result[i] = i
			}
			return result
		},
	}).ParseFiles("./frontend/templates/table.html"),
)

type Column struct {
	ID    string
	Label string
	Unit  string
	Cells []Cell
}

var table = []Column{
	{
		ID:    "a",
		Label: "a",
		Cells: []Cell{
			{Name: "a0"},
		},
	},
	{
		ID:    "b",
		Label: "b",
		Cells: []Cell{
			{Name: "b0"},
		},
	},
}

type Cell struct {
	Name  string
	Value string
}

func RowCount() int {
	if len(table) == 0 || len(table[0].Cells) == 0 {
		return 0
	}
	return len(table[0].Cells)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	Template.Execute(w, table)
}

func AddRowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, RedirectTO, http.StatusSeeOther)
		return
	}
	r.ParseForm()

	rowCount := RowCount()
	if rowCount == 0 {
		for i := range table {
			table[i].Cells = append(table[i].Cells, Cell{
				Name: fmt.Sprintf("%s0", table[i].ID),
			})
		}
		Template.ExecuteTemplate(w, "table_body_wrapper", table)
		return
	}

	lastIndex := rowCount - 1

	allFilled := true
	for i := range table {
		cell := &table[i].Cells[lastIndex]
		cell.Value = r.FormValue(cell.Name)
		if strings.TrimSpace(cell.Value) == "" {
			allFilled = false
		}
	}

	if !allFilled {
		Template.ExecuteTemplate(w, "table_body_wrapper", table)
		return
	}

	newIndex := rowCount
	for i := range table {
		table[i].Cells = append(table[i].Cells, Cell{
			Name: fmt.Sprintf("%s%d", table[i].ID, newIndex),
		})
	}

	Template.ExecuteTemplate(w, "table_body_wrapper", table)
}

func columnName(n int) string {
	result := ""
	for n >= 0 {
		result = string('a'+(n%26)) + result
		n = (n / 26) - 1
	}
	return result
}

func AddColumnHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Query := r.URL.Query()
	nextID := columnName(len(table))
	label := Query.Get("label")
	unit := Query.Get("unit")

	if label == "" {
		label = nextID
	}

	rowCount := RowCount()
	cells := make([]Cell, rowCount)
	for i := 0; i < rowCount; i++ {
		cells[i] = Cell{Name: fmt.Sprintf("%s%d", nextID, i)}
	}

	table = append(table, Column{
		ID:    nextID,
		Label: label,
		Unit:  unit,
		Cells: cells,
	})

	http.Redirect(w, r, RedirectTO, http.StatusSeeOther)
}

func RenameColumnHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Query := r.URL.Query()

	id := Query.Get("id")
	label := Query.Get("label")

	if id == "" || label == "" {
		http.Error(w, "missing id or label", http.StatusBadRequest)
		return
	}

	for i := range table {
		if table[i].ID == id {
			table[i].Label = label
			break
		}
	}
	http.Redirect(w, r, RedirectTO, http.StatusSeeOther)
}

func SavingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	r.ParseForm()

	name := r.FormValue("name")
	value := r.FormValue(name)

	for cell := range table {
		for row := range table[cell].Cells {
			cell := &table[cell].Cells[row]
			if cell.Name == name {
				cell.Value = value
				break
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewTableHandler(w http.ResponseWriter, r *http.Request) {
	table = []Column{
		{
			ID:    "a",
			Label: "a",
			Cells: []Cell{
				{Name: "a0"},
			},
		},
		{
			ID:    "b",
			Label: "b",
			Cells: []Cell{
				{Name: "b0"},
			},
		},
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.Redirect(w, r, RedirectTO, http.StatusSeeOther)
}
