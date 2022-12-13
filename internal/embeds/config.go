package embeds

import (
	_ "embed"
	"log"
	"os"
)

//go:embed template/sync-git.template.toml
var configTemplate string

// WriteConfigFile to disk
func WriteConfigFile() {
	if _, err := os.Stat("sync-git.toml"); os.IsNotExist(err) {
		err := os.WriteFile("sync-git.toml", []byte(configTemplate), 0644)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
	} else {
		log.Println("Config file already exists")
	}

}
