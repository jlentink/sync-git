package cnf

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

type GitDestination struct {
	Url      string
	Username string
	Password string
}

func (d *GitDestination) BaseName() string {
	pos := strings.LastIndex(d.Url, "/")
	if pos > -1 {
		return d.Url[pos+1:]
	}
	return d.Url
}

func (d *GitDestination) RemoteName() string {
	return strings.Replace(d.BaseName(), ".", "-", -1)
}

func (d *GitDestination) CloneUrl() string {
	matched, err := regexp.MatchString(`^http`, d.Url)
	if err != nil {
		log.Fatalf("Error matching regex %s", err)
	}
	if matched {
		u, err := url.Parse(d.Url)
		if err != nil {
			log.Fatalf("Error parsing url %s", err)
		}
		u.User = url.UserPassword(d.Username, d.Password)
		return u.String()
	}
	u := ""
	if d.Username != "" {
		u = d.Username
		if d.Password != "" {
			u += ":" + d.Password
		}
		u += "@"
	}
	u += d.Url
	return u
}
