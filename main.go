package main

import (
	"github.com/google/go-github/github"
	"fmt"
	"os"
	"os/exec"
)

var (
	sep = string(os.PathSeparator)
	path string
)

func main(){
	err := clone()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func clone() error {
	client := github.NewClient(nil)
	publicRepos, _, err := client.Repositories.ListByOrg("thisisfineio", nil)
	if err != nil {
		return err
	}
	gopath := os.Getenv("GOPATH")

	if path == "" {
		path = gopath + fmt.Sprintf("%ssrc%sgithub.com%sthisisfineio/test", sep, sep, sep)
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
				fmt.Println(string(e.Stderr))
			}
		} else {
			fmt.Print(string(o))
		}
	}
	return nil
}
