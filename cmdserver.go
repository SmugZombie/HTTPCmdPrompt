package main

// cmdserver.go
// Simple CMD Server for Windows
// Ron Egli - github.com/smugzombie
// Version 0.2

import (
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "os"
    "encoding/json"
)

var blacklist = ("shutdown")

// Design Command JSON Structure
type command_struct struct {
    Command string
}

func main() {
    // Serve Root
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, r.URL.Path[1:]) })
    // Test WHOAMI function
    http.HandleFunc("/whoami", whoami)
    // Kill the server
    http.HandleFunc("/kill", kill)
    // Run a command
    http.HandleFunc("/command", command)
    // Get CWD
    http.HandleFunc("/cwd", cwd)
    // Start the server
    log.Fatal(http.ListenAndServe(":8081", nil))
}

// Simply kill the server.
func kill(w http.ResponseWriter, req *http.Request){
    fmt.Fprintf(w, "Killing Server")
    _ = w
    _ = req
    os.Exit(0)
}

// Simply run WHOAMI
func whoami(w http.ResponseWriter, req *http.Request){
    fmt.Fprintf(w, runCommand( "whoami"))
}

// Get CWD
func cwd(w http.ResponseWriter, req *http.Request){
    fmt.Fprintf(w, runCommand("cd"))
    _ = w
    _ = req
}

// Run any command
func runCommand(command2run string) string{
    out, err := exec.Command("cmd","/c",command2run).Output()
    if err != nil { /*log.Println(err.Error())*/ return explainErrors(string(err.Error())) } 
    return string(out)
}

func explainErrors(errorcode string) string{
    if(errorcode == "exit status 1"){
        return "Invalid Command"
    }else if(errorcode == "exit status 2"){
        return "Invalid Permissions"
    }else{
        return "Undocumented Error!"
    }
}

// Function to call runCommand after parsing out command
func command(w http.ResponseWriter, req *http.Request) {
    // Expected get: ?{"command":"<enter command here>"}
    req.ParseForm()
    log.Println(req.Form)
    var t command_struct
    for key, _ := range req.Form {
        log.Println(key)
        err := json.Unmarshal([]byte(key), &t)
        if err != nil {
            log.Println(err.Error())
        }
    }
    fmt.Fprintf(w, runCommand(t.Command))
}
