build:
	go build -o dist/cmd cmd/cmd.go
run: build
	mkdir -p .ignore/content
	dist/cmd --table-id ${TABLE_ID} --token ${TOKEN} --export-path .ignore/content --image-path .ignore/images --template ./src/hugo.md.tpl --cmd hugo
