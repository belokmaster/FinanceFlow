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
	log.Printf("CreatePageHandler: Processing categories page request")

	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Printf("CreatePageHandler: Getting categories from database")
	categories, err := database.GetCategories(db)
	if err != nil {
		log.Printf("CreatePageHandler: Error getting categories: %v", err)
		http.Error(w, "could not get categories", http.StatusInternalServerError)
		return
	}

	log.Printf("CreatePageHandler: Retrieved %d accounts from database", len(categories))

	var categoriesForView []CategoryView
	for _, cat := range categories {
		color := "#" + cat.Color
		if color == "" {
			log.Printf("CreatePageHandler: Empty color for category %s, using default", cat.Name)
			color = "#4cd67a"
		}

		iconFileName, ok := database.IconCategoryFiles[cat.IconCode]
		if !ok {
			log.Printf("CreatePageHandler: Unknown icon code %d for category %s, using default icon", cat.IconCode, cat.Name)
			iconFileName = "Food"
		}

		iconHTML, ok := icons.CategoryIconCache[iconFileName]
		if !ok {
			log.Printf("CreatePageHandler: Icon file %s not found in cache for category %s, using food icon", iconFileName, cat.Name)
			iconHTML = icons.CategoryIconCache["Food"]
		}

		categoriesForView = append(categoriesForView, CategoryView{
			ID:       cat.ID,
			Name:     cat.Name,
			Color:    color,
			IconKey:  iconFileName,
			IconHTML: iconHTML,
		})
	}

	log.Printf("CreatePageHandler: Parsing template from path: %s", path)

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("CreatePageHandler: Error parsing template: %v", err)
		http.Error(w, "could not parse template", http.StatusInternalServerError)
		return
	}

	pageData := CategoryPageData{
		Categories: categoriesForView,
		Icons:      icons.CategoryIconCache,
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("CreatePageHandler: Error executing template: %v", err)
		log.Printf("error while proccess working: %v", err)
	} else {
		log.Printf("CreatePageHandler: Template executed successfully, sent %d categories to client", len(categoriesForView))
	}
}
