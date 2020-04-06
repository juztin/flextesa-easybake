package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
)

func handler(p *httputil.ReverseProxy, bake func()) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
		if r.Method != http.MethodGet && r.URL.Path == "/injection/operation" {
			bake()
		}
	}
}

func flextesa(args []string, s http.Server) (*exec.Cmd, func(), error) {
	args = append([]string{"mini-network"}, args...)
	pr, pw := io.Pipe()
	cmd := exec.Command("/usr/bin/flextesa", args...)

	fmt.Println(cmd.String())

	cmd.Stdin = pr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return cmd, nil, err
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Println(err)
			s.Shutdown(nil)
		}
	}()
	return cmd, func() {
		_, err := pw.Write([]byte("bake\n"))
		if err != nil {
			log.Println(err)
			s.Shutdown(nil)
		}
	}, err
}

func flextesaArgs(args []string) []string {
	for _, s := range args {
		if s == "--base-port" {
			log.Fatalln("Can't override Flextesa base-port")
		}
	}
	return append([]string{"--base-port", "30000"}, args...)
}

func main() {
	args := flextesaArgs(os.Args[1:])
	fmt.Println("Starting Flextesa/mini-network")
	mux := http.NewServeMux()
	s := http.Server{Addr: ":20000", Handler: mux}
	node, bake, err := flextesa(args, s)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Starting Easy-Bake...")
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   ":30000",
	})
	proxy.FlushInterval = -1
	mux.HandleFunc("/", handler(proxy, bake))
	err = s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Stopping...")
		node.Process.Kill()
		log.Fatalln(err)
	}
}
