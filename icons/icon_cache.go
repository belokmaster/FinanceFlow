package icons

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

//go:embed *
var iconsFS embed.FS

var IconCache = make(map[string]template.HTML)

var fileToName = map[string]string{
	"wallet": "Общий счет",
	"card":   "Карта",
	"coin":   "Наличные",
	"mark":   "Сберегательный счет",
	"house":  "Ипотека",
	"bag":    "Заем",
	"bonus":  "Бонусы",
}

func InitIcons() {
	icons, err := fs.ReadDir(iconsFS, ".")
	if err != nil {
		log.Fatalf("problem with open folder with icons: %v", err)
	}

	for _, iconFile := range icons {
		if iconFile.IsDir() {
			continue
		}

		if strings.HasSuffix(iconFile.Name(), ".go") {
			continue
		}

		content, err := fs.ReadFile(iconsFS, iconFile.Name())
		if err != nil {
			log.Printf("problem with reading icon %s: %v", iconFile.Name(), err)
			continue
		}
		contentStr := string(content)
		contentStr = strings.ReplaceAll(contentStr, `fill="#FFF"`, "")

		// delete .svg in end
		iconKey := strings.TrimSuffix(iconFile.Name(), filepath.Ext(iconFile.Name()))
		curName := fileToName[iconKey]
		IconCache[curName] = template.HTML(contentStr)
	}

	log.Printf("downloaded %d account's icons to cash", len(IconCache))
}
