prepare-cmd-test-wd:
	mkdir -p cmd/build

install: prepare-cmd-test-wd
	cd cmd && GOOS=darwin go build -o build/launcher_cmd

run: prepare-cmd-test-wd
	cd cmd/build && GOOS=darwin go run ../
