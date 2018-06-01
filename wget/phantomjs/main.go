package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func phantomjsGetURL(theURL string) (string, error) {
	const (
		cmdName    = "phantomjs"
		scriptName = "loadhtml.js"
	)
	var (
		cmdOut []byte
		err    error
	)

	parseBytes := func(buf []byte) string {
		serr := "TypeError: Attempting to change the setter of an unconfigurable property.\n"
		s := string(buf)
		s = strings.TrimSuffix(s, serr)
		s = strings.TrimSuffix(s, serr)
		return s
	}

	if cmdOut, err = exec.Command(cmdName, scriptName, theURL).Output(); err != nil {
		// exit status 1
		if strings.HasPrefix(err.Error(), "exit status ") {
			return "", fmt.Errorf("phantomjs: failed to load the URL %q", theURL)
		}
		return "", err
	}
	content := parseBytes(cmdOut)
	return content, nil
}

func main() {
	var (
		u    string
		html string
		err  error
	)

	u = "http://example.com"

	if html, err = phantomjsGetURL(u); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	fmt.Println(html)
}
