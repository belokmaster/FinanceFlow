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

	subcategoriesByParent := make(map[uint][]SubCategoryView)

	for _, subCat := range subCategories {
		color := "#" + subCat.Color
		if color == "" {
			log.Printf("CategoryPageHandler: Empty color for sub_category %s, using default", subCat.Name)
			color = "#4cd67a"
		}

		iconFileName, ok := database.IconSubCategoryFiles[subCat.IconCode]
		if !ok {
			log.Printf("CategoryPageHandler: Unknown icon code %d for sub_category %s, using default icon", subCat.IconCode, subCat.Name)
			iconFileName = "Restaraunt1"
		}

		iconHTML, ok := icons.SubCategoryIconCache[iconFileName]
		if !ok {
			log.Printf("CategoryPageHandler: Icon file %s not found in cache for sub_category %s, using food icon", iconFileName, subCat.Name)
			iconHTML = icons.SubCategoryIconCache["Restaraunt1"]
		}

		subcategoryView := SubCategoryView{
			ID:       subCat.ID,
			Name:     subCat.Name,
			Color:    color,
			IconKey:  iconFileName,
			IconHTML: iconHTML,
			ParentID: subCat.CategoryID,
		}

		subcategoriesByParent[subCat.CategoryID] = append(subcategoriesByParent[subCat.CategoryID], subcategoryView)
	}

	var categoriesForView []CategoryView
	for _, cat := range categories {
		color := "#" + cat.Color
		if color == "" {
			log.Printf("CategoryPageHandler: Empty color for category %s, using default", cat.Name)
			color = "#4cd67a"
		}

		iconFileName, ok := database.IconCategoryFiles[cat.IconCode]
		if !ok {
			log.Printf("CategoryPageHandler: Unknown icon code %d for category %s, using default icon", cat.IconCode, cat.Name)
			iconFileName = "Food"
		}

		iconHTML, ok := icons.CategoryIconCache[iconFileName]
		if !ok {
			log.Printf("CategoryPageHandler: Icon file %s not found in cache for category %s, using food icon", iconFileName, cat.Name)
			iconHTML = icons.CategoryIconCache["Food"]
		}

		subcats := subcategoriesByParent[cat.ID]
		if subcats == nil {
			subcats = []SubCategoryView{}
		}

		categoriesForView = append(categoriesForView, CategoryView{
			ID:            cat.ID,
			Name:          cat.Name,
			Color:         color,
			IconKey:       iconFileName,
			IconHTML:      iconHTML,
			Subcategories: subcats,
		})
	}

	log.Printf("CategoryPageHandler: Parsing template from path: %s", path)

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
		log.Printf("error while proccess working: %v", err)
	} else {
		log.Printf("CategoryPageHandler: Template executed successfully, sent %d categories to client", len(categoriesForView))
	}
}
