package http

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var hello = "hello world!"

func TestRun(t *testing.T) {
	conf := DefaultConfig
	conf.Addr = ":12345"

	handler := gin.Default()
	handler.GET("/", func(context *gin.Context) {
		fmt.Fprintln(context.Writer, hello)
	})

	if err := Run(handler, conf); err != nil {
		t.Fatal(err)
	}

	rep, err := http.Get("http://127.0.0.1:12345")
	if err != nil {
		t.Fatalf("get err:%s\n", err)
	}
	defer rep.Body.Close()

	if rep.StatusCode != http.StatusOK {
		t.Fatal("http status code invalid", rep.StatusCode)
	}

	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}

	result := strings.TrimRight(string(body), "\r\n")
	if result != hello {
		t.Fatalf("want [%s] real [%s]", hello, result)
	}
}

func TestRunTLS(t *testing.T) {
	conf := DefaultConfig
	conf.HTTPS = ":12346"

	handler := gin.Default()
	handler.GET("/", func(context *gin.Context) {
		fmt.Fprintln(context.Writer, hello)
	})

	if err := RunTLS(handler, conf); err != nil {
		t.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	rep, err := client.Get("https://127.0.0.1:12346")
	if err != nil {
		t.Fatalf("get err:%s\n", err)
	}
	defer rep.Body.Close()

	if rep.StatusCode != http.StatusOK {
		t.Fatal("http status code invalid", rep.StatusCode)
	}

	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}

	result := strings.TrimRight(string(body), "\r\n")
	if result != hello {
		t.Fatalf("want [%s] real [%s]", hello, result)
	}
}
