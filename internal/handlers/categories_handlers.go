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
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("CreateCategoryHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("CreateCategoryHandler: Form data received - %+v", r.Form)

	name := r.FormValue("Name")
	name = strings.TrimSpace(name)

	color := r.FormValue("Color")
	log.Printf("CreateCategoryHandler: Color '%s'", color)
	color = strings.TrimPrefix(color, "#")

	str_icon_id := r.FormValue("Icon")
	icon_id := database.IconCategoryNamesToIDs[str_icon_id]

	category := database.Category{
		Name:     name,
		Color:    color,
		IconCode: database.TypeCategoryIcons(icon_id),
	}

	log.Printf("CreateCategoryHandler: Creating category  Name=%s, Color=%s, Icon=%d", name, color, icon_id)

	err = database.AddCategory(db, category)
	if err != nil {
		log.Printf("CreateCategoryHandler: Failed to create category: %v", err)
		http.Error(w, fmt.Sprintf("failed to create a new category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateCategoryHandler: Category created successfully")
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("DeleteCategoryHandler: Starting category deletion")

	if r.Method != http.MethodPost {
		log.Printf("DeleteCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
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
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		log.Printf("UpdateCategoryHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("UpdateCategoryHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("UpdateCategoryHandler: Form data received - %+v", r.Form)

	categoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("UpdateCategoryHandler: Invalid account ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	newColor := r.FormValue("Color")
	if newColor != "" {
		err = database.ChangeCategoryColor(db, categoryID, newColor)
		if err != nil {
			log.Printf("UpdateCategoryHandler: Changing color to '%s' for category %d", newColor, categoryID)
			http.Error(w, fmt.Sprintf("failed to change category color: %v", err), http.StatusInternalServerError)
			return
		}
	}

	newName := r.FormValue("Name")
	if newName != "" {
		err = database.ChangeCategoryName(db, categoryID, newName)
		if err != nil {
			log.Printf("UpdateCategoryHandler: Changing name to '%s' for category %d", newName, categoryID)
			http.Error(w, fmt.Sprintf("failed to change category name: %v", err), http.StatusInternalServerError)
			return
		}
	}

	icon := r.FormValue("Icon")
	if icon != "" {
		iconCode, exists := database.IconCategoryNamesToIDs[icon]
		if !exists {
			log.Printf("UpdateCategoryHandler: Unknown icon name '%s'", icon)
			http.Error(w, "unknown icon name", http.StatusBadRequest)
			return
		}

		log.Printf("UpdateCategoryHandler: Changing icon for category ID %d to icon ID %d", categoryID, int(iconCode))

		err = database.ChangeCategoryIcon(db, categoryID, int(iconCode))
		if err != nil {
			log.Printf("UpdateCategoryHandler: Failed to change icon for category ID %d: %v", categoryID, err)
			http.Error(w, fmt.Sprintf("failed to change category's icon: %v", err), http.StatusInternalServerError)
			return
		}
	}

	log.Printf("UpdateAccountHandler: Successfully processed all changes for category %d", categoryID)
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
