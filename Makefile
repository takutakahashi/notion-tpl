build:
	go build -o dist/cmd cmd/cmd.go
run: build
	mkdir -p ./content
	dist/cmd --table-id ${TABLE_ID} --token ${NOTION_TOKEN} --export-path ./content --template ./src/hugo.md.tpl

