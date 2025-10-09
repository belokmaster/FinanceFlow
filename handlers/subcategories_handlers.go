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

func CreateSubCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	subCategoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if subCategoryID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
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

	subCategory := database.SubCategory{
		ID:         uint(subCategoryID),
		Name:       name,
		CategoryID: uint(categoryID),
	}

	err = database.AddSubCategory(db, subCategory)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create a new category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("A new sub_category was created sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteSubCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	subCategoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if subCategoryID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	err = database.DeleteSubCategory(db, subCategoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete sub_category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("sub_category was deleted sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
