build:
	go build -o goshell ./cmd/goshell
	chmod +x goshell

test-unit: build
	go test ./tests/unit -v

test-integration: build
	go test ./tests/integration -v

test: test-unit test-integration

clean:
	rm -f goshell
	rm -f testfile.txt

gits_up:
	git status
	git add .
	git commit -m "Update"
	git push
