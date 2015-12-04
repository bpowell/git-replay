git remote add temp file:///path/to/git/repo

echo "[" > data.json; git log temp/master --pretty=oneline --reverse | awk '{print "{\"commit\":\"" $1 "\",\"msg\":\"" substr($0, 42) "\"},"}' >> data.json; sed -i '$ s/.$//' data.json; echo "]" >> data.json

copy main.go to repo

go run main.go
