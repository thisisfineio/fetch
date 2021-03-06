package main

import (
	"github.com/google/go-github/github"
	"fmt"
	"os"
	"os/exec"
	"github.com/thisisfineio/variant"
	"time"
	"flag"
	"strings"
)

var (
	sep = string(os.PathSeparator)
	path string
	gopath = os.Getenv("GOPATH")
	org string
	update bool
	// inline versions for now, maybe future use build tool to code gen and read from json
	cur = &variant.Version{
		Major: 0,
		Minor: 1,
		Date: time.Unix(1482639633, 0),
		ReleaseType: variant.Beta,
		Description: "The initial beta releae of 'fetch', designed to make it easy for developers to grab all of thisisfine.io's code and start contributing.",
	}

	Versions = &variant.Versions{
		Current: cur,
		History: []*variant.Version{
			cur,
		},
	}

	v = flag.Bool("version", false, "Prints the current version of fetch")
)

func init(){
	flag.StringVar(&path, "path", gopath + fmt.Sprintf("%ssrc%sgithub.com%sthisisfineio", sep, sep, sep), "Sets the path to download repos to")
	flag.StringVar(&org, "org", "thisisfineio", "Sets the organization to clone from")
	flag.BoolVar(&update, "u", false, "Attempt to update existing repositories - will not stash changes, will fail if unstaged things exist")
}

func main(){
	parse()
	err := clone()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func clone() error {
	client := github.NewClient(nil)
	publicRepos, _, err := client.Repositories.ListByOrg(org, nil)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}

	err = os.Chdir(path)
	if err != nil {
		return err
	}

	cmds := make([]*exec.Cmd, 0)
	for _, r := range publicRepos {
		cmd := exec.Command("git", "clone", *r.CloneURL)
		cmd.Env = os.Environ()
		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		o, err := cmd.Output()
		if err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				if strings.Contains(string(e.Stderr), "already exists and is not an empty directory.") {
					// put git pull here
					continue
				}
				fmt.Println(string(e.Stderr))
			} else {
				fmt.Print(string(o))
			}
		}
	}
	return nil
}

func parse() {
	flag.Parse()
	if *v {
		version()
	}
}

func version() {
	fmt.Println("fetch version:", cur.VersionString())
	os.Exit(0)
}
