package handlers

import (
	"finance_flow/internal/database"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("CreateCategoryHandler: Starting category creation")

	if r.Method != http.MethodPost {
		log.Printf("CreateCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("CreateCategoryHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("CreateCategoryHandler: Form data received - %+v", r.Form)

	categoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("CreateCategoryHandler: Invalid category ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		log.Printf("CreateCategoryHandler: Negative category ID %d", categoryID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	name := r.FormValue("Name")
	name = strings.TrimSpace(name)
	log.Printf("CreateCategoryHandler: Category name '%s' (trimmed)", name)

	category := database.Category{
		ID:   uint(categoryID),
		Name: name,
	}

	log.Printf("CreateCategoryHandler: Creating category with ID=%d, Name=%s", categoryID, name)

	err = database.AddCategory(db, category)
	if err != nil {
		log.Printf("CreateCategoryHandler: Failed to create category: %v", err)
		http.Error(w, fmt.Sprintf("failed to create a new category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateCategoryHandler: Category created successfully - ID=%d, Name=%s", categoryID, name)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("DeleteCategoryHandler: Starting category deletion")

	if r.Method != http.MethodPost {
		log.Printf("DeleteCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("DeleteCategoryHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteCategoryHandler: Form data received - %+v", r.Form)

	categoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("DeleteCategoryHandler: Invalid category ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		log.Printf("DeleteCategoryHandler: Negative category ID %d", categoryID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteCategoryHandler: Attempting to delete category ID %d", categoryID)

	err = database.DeleteCategory(db, categoryID)
	if err != nil {
		log.Printf("DeleteCategoryHandler: Failed to delete category ID %d: %v", categoryID, err)
		http.Error(w, fmt.Sprintf("failed to delete category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteCategoryHandler: Category ID %d deleted successfully", categoryID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
