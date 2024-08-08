package main

import (
	"bytes"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"strings"
)

func pullSeeds(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return []string{}, err
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 1000000 /*1MB*/))
	if err != nil {
		return []string{}, err
	}
	if len(body) == 0 {
		return []string{}, err
	}
	lines := bytes.Split(bytes.ReplaceAll(body, []byte("\r\n"), []byte("\n")), []byte("\n"))
	result := make([]string, 0, len(lines))
	for i := range lines {
		line := strings.TrimSpace(string(lines[i]))
		if len(line) == 0 {
			continue
		}
		if slices.Contains(result, line) {
			continue
		}
		result = append(result, line)
	}
	return result, nil
}
