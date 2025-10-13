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

func CreateSubCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("CreateSubCategoryHandler: Starting subcategory creation")

	if r.Method != http.MethodPost {
		log.Printf("CreateSubCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("CreateSubCategoryHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("CreateSubCategoryHandler: Form data received - %+v", r.Form)

	name := r.FormValue("Name")
	name = strings.TrimSpace(name)
	log.Printf("CreateSubCategoryHandler: Subcategory name '%s'", name)

	color := r.FormValue("Color")
	log.Printf("CreateSubCategoryHandler: Color '%s'", color)
	color = strings.TrimPrefix(color, "#")

	str_icon_id := r.FormValue("Icon")
	icon_id := database.IconSubCategoryNamesToIDs[str_icon_id]

	parentIDStr := r.FormValue("ParentID")
	if parentIDStr == "" {
		log.Printf("CreateSubCategoryHandler: Error - ParentID is missing from the form")
		http.Error(w, "ParentID is required", http.StatusBadRequest)
		return
	}

	parentID_64, err := strconv.ParseUint(parentIDStr, 10, 32)
	if err != nil {
		log.Printf("CreateSubCategoryHandler: Error parsing ParentID '%s': %v", parentIDStr, err)
		http.Error(w, "Invalid ParentID", http.StatusBadRequest)
		return
	}
	parentID := uint(parentID_64)

	subCategory := database.SubCategory{
		Name:       name,
		Color:      color,
		IconCode:   database.TypeSubCategoryIcons(icon_id),
		CategoryID: parentID,
	}

	log.Printf("CreateSubCategoryHandler: Creating sub_category  Name=%s, Color=%s, Icon=%d", name, color, icon_id)

	err = database.AddSubCategory(db, subCategory)
	if err != nil {
		log.Printf("CreateSubCategoryHandler: Failed to create subcategory: %v", err)
		http.Error(w, fmt.Sprintf("failed to create a new category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateSubCategoryHandler: Subcategory created successfully.")
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func DeleteSubCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("DeleteSubCategoryHandler: Starting subcategory deletion")

	if r.Method != http.MethodPost {
		log.Printf("DeleteSubCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
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
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func UpdateSubCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		log.Printf("UpdateSubCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("UpdateSubCategoryHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("UpdateSubCategoryHandler: Form data received - %+v", r.Form)

	subCategoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("UpdateSubCategoryHandler: Invalid account ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	newColor := r.FormValue("Color")
	if newColor != "" {
		err = database.ChangeSubCategoryColor(db, subCategoryID, newColor)
		if err != nil {
			log.Printf("UpdateSubCategoryHandler: Changing color to '%s' for sub_category %d", newColor, subCategoryID)
			http.Error(w, fmt.Sprintf("failed to change sub_category color: %v", err), http.StatusInternalServerError)
			return
		}
	}

	newName := r.FormValue("Name")
	if newName != "" {
		err = database.ChangeSubCategoryName(db, subCategoryID, newName)
		if err != nil {
			log.Printf("UpdateSubCategoryHandler: Changing name to '%s' for sub_category %d", newName, subCategoryID)
			http.Error(w, fmt.Sprintf("failed to change sub_category name: %v", err), http.StatusInternalServerError)
			return
		}
	}

	icon := r.FormValue("Icon")
	if icon != "" {
		iconCode, exists := database.IconSubCategoryNamesToIDs[icon]
		if !exists {
			log.Printf("UpdateSubCategoryHandler: Unknown icon name '%s'", icon)
			http.Error(w, "unknown icon name", http.StatusBadRequest)
			return
		}

		log.Printf("UpdateSubCategoryHandler: Changing icon for sub_category ID %d to icon ID %d", subCategoryID, int(iconCode))

		err = database.ChangeSubCategoryIcon(db, subCategoryID, int(iconCode))
		if err != nil {
			log.Printf("UpdateSubCategoryHandler: Failed to change icon for sub_category ID %d: %v", subCategoryID, err)
			http.Error(w, fmt.Sprintf("failed to change sub_category's icon: %v", err), http.StatusInternalServerError)
			return
		}
	}

	log.Printf("UpdateSubCategoryHandler: Successfully processed all changes for category %d", subCategoryID)
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
