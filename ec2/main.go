package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func assert(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func get(dir []string) io.Reader {
	url := filepath.Join(dir...)
	req, err := http.Get(fmt.Sprintf("http://169.254.169.254/latest/meta-data/%s", url))
	assert(err)
	return req.Body
}

func getAll(dir []string) string {
	data, err := ioutil.ReadAll(get(dir))
	assert(err)
	return string(data)
}

func getDir(acc map[string]interface{}, dir []string) {
	s := bufio.NewScanner(get(dir))
	for s.Scan() {
		line := s.Text()
		if line[len(line)-1] == '/' {
			part := line[0 : len(line)-1]
			acc[part] = make(map[string]interface{})
			getDir(acc[part].(map[string]interface{}), append(dir, part))
		} else if len(dir) > 0 && dir[len(dir)-1] == "public-keys" {
			parts := strings.Split(line, "=")
			acc[parts[1]] = getAll(append(dir, []string{parts[0], "openssh-key"}...))
		} else {
			acc[line] = getAll(append(dir, line))
		}
	}
	assert(s.Err())
}

func main() {
	m := make(map[string]interface{})
	getDir(m, []string{})
	json, err := json.Marshal(m)
	assert(err)
	fmt.Println(string(json))
}
