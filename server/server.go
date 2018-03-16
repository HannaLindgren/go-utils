package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func getParam(paramName string, r *http.Request) string {
	res := r.FormValue(paramName)
	if res != "" {
		return res
	}
	res = r.PostFormValue(paramName)
	if res != "" {
		return res
	}
	vars := mux.Vars(r)
	return vars[paramName]
}

func execCmd(cmd *exec.Cmd) (bytes.Buffer, bytes.Buffer, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("command: %v", strings.Join(cmd.Args, " "))

	return stdout, stderr, cmd.Run()
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	cmdArgsWithVars := []string{}
	for _, s := range cmdArgs {
		cmdArgsWithVars = append(cmdArgsWithVars, s)
	}
	for _, name := range cmdVars {
		value := getParam(name, r)
		if value != "" {
			name = fmt.Sprintf("{%s}", name)
			for i, s := range cmdArgsWithVars {
				cmdArgsWithVars[i] = strings.Replace(s, name, value, -1)
			}
		}
	}

	cmd := exec.Command(cmdName, cmdArgsWithVars...)
	stdout, stderr, err := execCmd(cmd)
	result := strings.TrimSpace(stdout.String())
	if result == "" {
		result = "<empty output>"
	}
	stderrString := strings.TrimSpace(stderr.String())
	log.Printf("result: %s\n", result)
	if err != nil {
		log.Printf("error: %v\n", err)
		if stderrString == "" {
			stderrString = "<empty>"
		}
		log.Printf("stderr: %s\n", stderrString)
		msg := fmt.Sprintf("failed running '%s': %v\n", cmd.Path, err)
		log.Print(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	// if no error:
	if stderrString != "" {
		log.Printf("stderr: %s\n", stderrString)
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "result: %s\n", result)
}

var cmdVarRe = regexp.MustCompile("{([^}]+)}")

func parseCmdVars(cmd string) []string {
	res := []string{}
	matches := cmdVarRe.FindAllStringSubmatch(cmd, -1)
	for _, m := range matches {
		varName := m[1]
		res = append(cmdVars, varName)
	}
	return res
}

var cmdName string
var cmdArgs []string
var cmdVars []string

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "server 'PORT' 'COMMAND'\n")
		os.Exit(0)
	}

	port := os.Args[1]
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	cmd := os.Args[2]
	cmdSplit := strings.Fields(cmd)
	cmdName = cmdSplit[0]
	cmdArgs = cmdSplit[1:]
	cmdVars = parseCmdVars(cmd) // []string{}

	url := "/"
	// for PARAMS
	// for _, v := range cmdVars {
	// 	url = fmt.Sprintf("%s{%s}/", url, v)
	// }

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc(url, handlerFunc)

	log.Printf("starting g2p server at port: %s\n", port)
	log.Printf("responding to url: %s", url)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("no fun: %v\n", err)
	}

}
