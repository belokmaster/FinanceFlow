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

func InitIcons() {
	icons, err := fs.ReadDir(iconsFS, ".")
	if err != nil {
		log.Fatalf("problem with open folder with icons: %v", err)
	}

	for _, iconFile := range icons {
		if iconFile.IsDir() {
			continue
		}

		content, err := fs.ReadFile(iconsFS, iconFile.Name())
		if err != nil {
			log.Printf("problem with reading icon %s: %v", iconFile.Name(), err)
			continue
		}
		// delete .svg in end
		iconKey := strings.TrimSuffix(iconFile.Name(), filepath.Ext(iconFile.Name()))
		IconCache[iconKey] = template.HTML(content)
	}
	log.Printf("downloaded %d icons to cash", len(IconCache))
}
