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
	log.Printf("CreateSubCategoryHandler: Starting subcategory creation")

	if r.Method != http.MethodPost {
		log.Printf("CreateSubCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("CreateSubCategoryHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("CreateSubCategoryHandler: Form data received - %+v", r.Form)

	subCategoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("CreateSubCategoryHandler: Invalid subcategory ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if subCategoryID < 0 {
		log.Printf("CreateSubCategoryHandler: Negative subcategory ID %d", subCategoryID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("CategoryID"))
	if err != nil {
		log.Printf("CreateSubCategoryHandler: Invalid category ID '%s': %v", r.FormValue("CategoryID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		log.Printf("CreateSubCategoryHandler: Negative category ID %d", categoryID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	name := r.FormValue("Name")
	name = strings.TrimSpace(name)
	log.Printf("CreateSubCategoryHandler: Subcategory name '%s'", name)

	subCategory := database.SubCategory{
		ID:         uint(subCategoryID),
		Name:       name,
		CategoryID: uint(categoryID),
	}

	log.Printf("CreateSubCategoryHandler: Creating subcategory with ID=%d, Name=%s, Parent Category ID=%d", subCategoryID, name, categoryID)

	err = database.AddSubCategory(db, subCategory)
	if err != nil {
		log.Printf("CreateSubCategoryHandler: Failed to create subcategory: %v", err)
		http.Error(w, fmt.Sprintf("failed to create a new category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateSubCategoryHandler: Subcategory created successfully - ID=%d, Name=%s, Parent Category ID=%d", subCategoryID, name, categoryID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteSubCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("DeleteSubCategoryHandler: Starting subcategory deletion")

	if r.Method != http.MethodPost {
		log.Printf("DeleteSubCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("DeleteSubCategoryHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteSubCategoryHandler: Form data received - %+v", r.Form)

	subCategoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("DeleteSubCategoryHandler: Invalid subcategory ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if subCategoryID < 0 {
		log.Printf("DeleteSubCategoryHandler: Negative subcategory ID %d", subCategoryID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteSubCategoryHandler: Attempting to delete subcategory ID %d", subCategoryID)

	err = database.DeleteSubCategory(db, subCategoryID)
	if err != nil {
		log.Printf("DeleteSubCategoryHandler: Failed to delete subcategory ID %d: %v", subCategoryID, err)
		http.Error(w, fmt.Sprintf("failed to delete sub_category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteSubCategoryHandler: Subcategory ID %d deleted successfully", subCategoryID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
