package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Hash struct {
	Commit string
	Msg    string
}

type Branch struct {
	Name    string
	Commits []Hash
}

func (h Hash) commit() {
	var out bytes.Buffer
	fmt.Println(h.Msg)
	args := strings.Split("cherry-pick --strategy-option theirs", " ")
	args = append(args, h.Commit)
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error with cherry-pick")
		fmt.Println(err)
	}

	fmt.Println(out.String())
}

func pushMaster() {
	fmt.Println("pushing to master")
	args := strings.Split("push", " ")
	cmd := exec.Command("git", args...)
	fmt.Println(cmd.Run())
}

func pullMaster() {
	fmt.Println("pulling master")
	args := strings.Split("pull", " ")
	cmd := exec.Command("git", args...)
	fmt.Println(cmd.Run())
}

func checkoutMaster() {
	fmt.Println("checking out master")
	args := strings.Split("checkout master", " ")
	cmd := exec.Command("git", args...)
	fmt.Println(cmd.Run())
}

func main() {
	file, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}

	var hashes []Hash
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&hashes)
	if err != nil {
		panic(err)
	}

	var branches []Branch

	hashes[0].commit()
	pushMaster()

	hashes = append(hashes[:0], hashes[1:]...)

	start := 0
	end := 0
	for _, hash := range hashes {
		if strings.Contains(hash.Msg, "Merge branch '") {
			var branch Branch
			branch.Name = strings.Split(hash.Msg, "'")[1]
			branch.Commits = append(branch.Commits, hashes[start:end]...)
			branches = append(branches, branch)

			start = end + 1
		}

		end = end + 1
	}

	for _, b := range branches {
		fmt.Println(b)
		fmt.Println("")

		var out bytes.Buffer
		args := strings.Split("checkout -b", " ")
		args = append(args, b.Name)
		cmd := exec.Command("git", args...)
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error with checkout -b")
			fmt.Println(err)
		}
		fmt.Println(out.String())

		for _, h := range b.Commits {
			h.commit()
		}

		args = strings.Split("push --set-upstream origin", " ")
		args = append(args, b.Name)
		cmd = exec.Command("git", args...)
		fmt.Println(cmd.Run())

		var nilString string
		fmt.Println("Waiting for merge")
		fmt.Scanf("%s", &nilString)

		checkoutMaster()
		pullMaster()
	}
}
