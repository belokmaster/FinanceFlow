package icons

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

//go:embed account_icons/*.svg
var accountIconsFS embed.FS

//go:embed category_icons/*.svg
var categoryIconsFS embed.FS

//go:embed subcategory_icons/*.svg
var subCategoryIconsFS embed.FS

func InitIcons() {
	initAccountIcons()
	initCategoryIcons()
	initSubCategoryIcons()
	log.Printf("Initialized %d account icons, %d category icons, %d sub_category icons",
		len(AccountIconCache), len(CategoryIconCache), len(SubCategoryIconCache))
}

func initAccountIcons() {
	icons, err := fs.ReadDir(accountIconsFS, "account_icons")
	if err != nil {
		log.Printf("Error loading account icons: %v", err)
		return
	}

	for _, iconFile := range icons {
		if iconFile.IsDir() {
			continue
		}

		filePath := filepath.Join("account_icons", iconFile.Name())
		content, err := fs.ReadFile(accountIconsFS, filePath)
		if err != nil {
			log.Printf("Error reading icon %s: %v", iconFile.Name(), err)
			continue
		}

		contentStr := string(content)

		iconKey := strings.TrimSuffix(iconFile.Name(), filepath.Ext(iconFile.Name()))

		curName, ok := accountFileToName[iconKey]
		if !ok {
			log.Printf("Warning: no display name for account icon %s", iconFile.Name())
			continue
		}
		AccountIconCache[curName] = template.HTML(contentStr)
	}
}

func initCategoryIcons() {
	icons, err := fs.ReadDir(categoryIconsFS, "category_icons")
	if err != nil {
		log.Printf("Error loading category icons: %v", err)
		return
	}

	for _, iconFile := range icons {
		if iconFile.IsDir() {
			continue
		}

		filePath := filepath.Join("category_icons", iconFile.Name())
		content, err := fs.ReadFile(categoryIconsFS, filePath)
		if err != nil {
			log.Printf("Error reading icon %s: %v", iconFile.Name(), err)
			continue
		}

		contentStr := string(content)
		iconKey := strings.TrimSuffix(iconFile.Name(), filepath.Ext(iconFile.Name()))

		curName, ok := categoryFileToName[iconKey]
		if !ok {
			log.Printf("Warning: no display name for category icon %s", iconFile.Name())
			continue
		}
		CategoryIconCache[curName] = template.HTML(contentStr)
	}
}

func initSubCategoryIcons() {
	icons, err := fs.ReadDir(subCategoryIconsFS, "subcategory_icons")
	if err != nil {
		log.Printf("Error loading sub_category icons: %v", err)
		return
	}

	for _, iconFile := range icons {
		if iconFile.IsDir() {
			continue
		}

		filePath := filepath.Join("subcategory_icons", iconFile.Name())
		content, err := fs.ReadFile(subCategoryIconsFS, filePath)
		if err != nil {
			log.Printf("Error reading icon %s: %v", iconFile.Name(), err)
			continue
		}

		contentStr := string(content)
		iconKey := strings.TrimSuffix(iconFile.Name(), filepath.Ext(iconFile.Name()))

		curName, ok := subCategoryFileToName[iconKey]
		if !ok {
			log.Printf("Warning: no display name for sub_category icon %s", iconFile.Name())
			continue
		}
		SubCategoryIconCache[curName] = template.HTML(contentStr)
	}
}

// getters
func GetAccountIcon(name string) template.HTML {
	icon, ok := AccountIconCache[name]
	if !ok {
		log.Printf("account icon not found for %s", name)
		return template.HTML("")
	}
	return icon
}

func GetCategoryIcon(name string) template.HTML {
	icon, ok := CategoryIconCache[name]
	if !ok {
		log.Printf("category icon not found for %s", name)
		return template.HTML("")
	}
	return icon
}

func GetSubCategoryIcon(name string) template.HTML {
	icon, ok := SubCategoryIconCache[name]
	if !ok {
		log.Printf("sub_category icon not found for %s", name)
		return template.HTML("")
	}
	return icon
}

func GetAccountIcons() map[string]template.HTML {
	return AccountIconCache
}

func GetCategoryIcons() map[string]template.HTML {
	return CategoryIconCache
}

func GetSubCategoryIcons() map[string]template.HTML {
	return SubCategoryIconCache
}
