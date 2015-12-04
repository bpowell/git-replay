git remote add temp file:///path/to/git/repo

git log temp/master --pretty=oneline --reverse | awk '{print "{\"commit\":\"" $1 "\",\"msg\":\"" substr($0, 42) "\"},"}' > data.json

copy main.go to repo

go run main.go
