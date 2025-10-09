package handlers

import (
	"finance_flow/database"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	name := r.FormValue("Name")
	name = strings.TrimSpace(name)

	category := database.Category{
		ID:   uint(categoryID),
		Name: name,
	}

	err = database.AddCategory(db, category)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create a new category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("A new category was created sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	err = database.DeleteCategory(db, categoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("category was deleted sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
