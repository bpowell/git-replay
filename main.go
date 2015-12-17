package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	Id      int `json: "id"`
	Commits []Hash
}

func (b *Branch) submitMergeRequest() {
	client := &http.Client{}
	values := url.Values{}
	values.Add("id", "349")
	values.Add("source_branch", b.Name)
	values.Add("target_branch", "master")
	values.Add("title", b.Name)

	req, _ := http.NewRequest("POST", "https://code.com/api/v3/projects/349/merge_requests", strings.NewReader(values.Encode()))
	req.Header.Add("PRIVATE-TOKEN", "")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))

	_ = json.Unmarshal(body, b)
	fmt.Println("submitted merge request")
}

func (b *Branch) acceptMergeRequest() {
	client := &http.Client{}
	values := url.Values{}
	values.Add("id", "349")
	values.Add("merge_request_id", fmt.Sprintf("%d", b.Id))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("https://code.com/api/v3/projects/349/merge_request/%d/merge", b.Id), strings.NewReader(values.Encode()))
	req.Header.Add("PRIVATE-TOKEN", "")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("acceptMergeRequest")
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

		fmt.Println("Waiting for merge")
		b.submitMergeRequest()
		b.acceptMergeRequest()

		checkoutMaster()
		pullMaster()
	}
}
