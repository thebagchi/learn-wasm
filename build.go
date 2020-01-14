package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CopyFile(src string, dst string) {
	data, err := ioutil.ReadFile(src)
	err = ioutil.WriteFile(dst, data, 0644)
	_ = err
}

func ParseEnv(cmd string) map[string]string {
	dump, err := exec.Command(cmd, strings.Fields("env")...).CombinedOutput()
	if nil != err {
		fmt.Println("Error: ", err)
		fmt.Println("Log: ", string(dump))
		return nil
	}
	env := map[string]string{
		// Empty
	}
	items := strings.Split(string(dump), "\n")
	for _, item := range items {
		item = strings.TrimSpace(item)
		items := strings.Split(item, "=")
		if len(items) == 2 {
			var (
				key   = items[0]
				value = items[1]
			)
			env[key] = strings.TrimPrefix(strings.TrimSuffix(value, "\""), "\"")
		}
	}
	return env
}

func BuildGo() {
	env := ParseEnv("go")
	if nil == env || len(env) == 0 {
		fmt.Println("Failed executing command \"go env\"")
		return
	}
	var (
		goos   = os.Getenv("GOOS")
		goarch = os.Getenv("GOARCH")
		root   = env["GOROOT"]
	)
	_ = os.Setenv("GOOS", "js")
	_ = os.Setenv("GOARCH", "wasm")
	defer func() {
		_ = os.Setenv("GOOS", goos)
		_ = os.Setenv("GOARCH", goarch)
	}()
	dump, err := exec.Command(
		"go",
		strings.Fields("build -o www/main.wasm wasm_go/main.go")...,
	).CombinedOutput()
	if nil != err {
		fmt.Println("Error: ", err)
		fmt.Println("Log: ", string(dump))
		return
	}
	CopyFile(filepath.Join(root, "/misc/wasm/wasm_exec.js"), "www/wasm_exec.js")
	CopyFile("wasm_go/index.html", "www/index.html")
}

func BuildRs() {
	dump, err := exec.Command(
		"wasm-pack",
		strings.Fields("build wasm_rs/")...,
	).CombinedOutput()
	if nil != err {
		fmt.Println("Error: ", err)
		fmt.Println("Log: ", string(dump))
		return
	}
	CopyFile("wasm_rs/index.html", "www/index.html")
}

func main() {
	var mode = flag.String("mode", "go", "oneof go/rs")
	flag.Parse()
	if *mode != "go" && *mode != "rs" {
		fmt.Println("Can only build for rust or golang")
		os.Exit(-1)
	}
	if *mode == "go" {
		BuildGo()
		return
	}
	if *mode == "rs" {
		BuildRs()
		return
	}
}
