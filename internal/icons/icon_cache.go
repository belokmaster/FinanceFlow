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

var accountFileToName = map[string]string{
	"wallet": "Общий счет",
	"card":   "Карта",
	"coin":   "Наличные",
	"mark":   "Сберегательный счет",
	"house":  "Ипотека",
	"bag":    "Заем",
	"bonus":  "Бонусы",
}

var categoryFileToName = map[string]string{
	"food": "Еда",
}

var AccountIconCache = make(map[string]template.HTML)
var CategoryIconCache = make(map[string]template.HTML)

func InitIcons() {
	initAccountIcons()
	initCategoryIcons()
	log.Printf("Initialized %d account icons and %d category icons",
		len(AccountIconCache), len(CategoryIconCache))
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
		contentStr = strings.ReplaceAll(contentStr, `fill="#FFF"`, "")

		iconKey := strings.TrimSuffix(iconFile.Name(), filepath.Ext(iconFile.Name()))
		curName := accountFileToName[iconKey]
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
		curName := categoryFileToName[iconKey]
		CategoryIconCache[curName] = template.HTML(contentStr)
	}
}

// getters
func GetAccountIcon(name string) template.HTML {
	return AccountIconCache[name]
}

func GetCategoryIcon(name string) template.HTML {
	return CategoryIconCache[name]
}

func GetAccountIcons() map[string]template.HTML {
	return AccountIconCache
}

func GetCategoryIcons() map[string]template.HTML {
	return CategoryIconCache
}
