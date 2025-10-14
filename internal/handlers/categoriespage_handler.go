package handlers

import (
	"finance_flow/internal/database"
	"finance_flow/internal/icons"
	"log"
	"net/http"
	"text/template"

	"gorm.io/gorm"
)

func CategoryPageHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB, path string) {
	log.Printf("CategoryPageHandler: Processing categories page request")

	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	log.Printf("CategoryPageHandler: Getting categories from database")
	categories, err := database.GetCategories(db)
	if err != nil {
		log.Printf("CreateCategoryPageHandler: Error getting categories: %v", err)
		http.Error(w, "could not get categories", http.StatusInternalServerError)
		return
	}

	log.Printf("CategoryPageHandler: Retrieved %d categories from database", len(categories))

	log.Printf("CategoryPageHandler: Getting categories from database")
	subCategories, err := database.GetSubCategories(db)
	if err != nil {
		log.Printf("CreateCategoryPageHandler: Error getting sub_categories: %v", err)
		http.Error(w, "could not get sub_categories", http.StatusInternalServerError)
		return
	}

	log.Printf("CategoryPageHandler: Retrieved %d sub_categories from database", len(categories))

	subcategoriesByParent := convertSubcategoriesToView(subCategories)
	categoriesForView := convertCategoriesToView(categories, subcategoriesByParent)

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("CategoryPageHandler: Error parsing template: %v", err)
		http.Error(w, "could not parse template", http.StatusInternalServerError)
		return
	}

	pageData := CategoryPageData{
		Categories:       categoriesForView,
		CategoryIcons:    icons.CategoryIconCache,
		SubcategoryIcons: icons.SubCategoryIconCache,
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("CategoryPageHandler: Error executing template: %v", err)
	} else {
		log.Printf("CategoryPageHandler: Template executed successfully, sent %d categories to client", len(categoriesForView))
	}
}
