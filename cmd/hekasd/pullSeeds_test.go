package main

import (
	"golang.org/x/exp/slices"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestPullSeeds(t *testing.T) {
	items := strings.Join([]string{
		"",
		"\ttest ",
		"\t\t  abcdefg  ",
		"",
		"foobar",
	}, "\r\n")
	shoudBe := []string{"test", "abcdefg", "foobar"}
	mux := http.NewServeMux()
	mux.HandleFunc("/seeds.txt", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writer.Write([]byte(items))
	})
	socket, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	go http.Serve(socket, mux)
	whatGot, err := pullSeeds((&url.URL{
		Scheme: "http",
		Host:   socket.Addr().String(),
		Path:   "/seeds.txt",
	}).String())
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(whatGot, shoudBe) {
		t.Fatalf("expected %q but got %q", shoudBe, whatGot)
	}
}
