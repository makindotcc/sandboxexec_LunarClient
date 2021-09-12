prepare-cmd-test-wd:
	mkdir -p cmd/build
	cp textures.txt cmd/build/textures.txt

install:
	$(prepare-cmd-test-wd)
	cd cmd && GOOS=darwin go build -o build/launcher_cmd

run:
	$(prepare-cmd-test-wd)
	cd cmd/build && GOOS=darwin go run ../
