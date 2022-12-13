package gitActions

import (
	"strings"
	"sync-git/internal/cnf"
)

func appendExtension(name string) string {
	if !strings.HasSuffix(name, ".git") {
		name += ".git"
	}
	return name
}

func getBaseName() string {
	url := cnf.GetString("source.url")
	pos := strings.LastIndex(url, "/")
	if pos > -1 {
		return appendExtension(url[pos+1:])
	}
	return appendExtension(url)
}
