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
	"food":        "Еда",
	"housing1":    "Жилье1",
	"housing2":    "Жилье2",
	"car_repair1": "Транспортное средство1",
	"car_repair2": "Транспортное средство2",
	"coin_cat":    "Доход",
	"computer1":   "Связь, ПК1",
	"computer2":   "Связь, ПК2",
	"bus1":        "Общественный транспорт1",
	"bus2":        "Общественный транспорт2",
	"community1":  "Жизнь и развлечения1",
	"community2":  "Жизнь и развлечения2",
	"community3":  "Жизнь и развлечения3",
	"expense1":    "Финансовые расходы1",
	"expense2":    "Финансовые расходы2",
}

var subCategoryFileToName = map[string]string{
	"restaurant1": "Ресторан1",
	"pizza1":      "Фаст-фуд1",
}

var AccountIconCache = make(map[string]template.HTML)
var CategoryIconCache = make(map[string]template.HTML)
var SubCategoryIconCache = make(map[string]template.HTML)

func InitIcons() {
	initAccountIcons()
	initCategoryIcons()
	initSubCategoryIcons()
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

func initSubCategoryIcons() {
	icons, err := fs.ReadDir(subCategoryIconsFS, "subcategory_icons")
	if err != nil {
		log.Printf("Error loading category icons: %v", err)
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
		curName := subCategoryFileToName[iconKey]
		SubCategoryIconCache[curName] = template.HTML(contentStr)
	}
}

// getters
func GetAccountIcon(name string) template.HTML {
	return AccountIconCache[name]
}

func GetCategoryIcon(name string) template.HTML {
	return CategoryIconCache[name]
}

func GetSubCategoryIcon(name string) template.HTML {
	return SubCategoryIconCache[name]
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
