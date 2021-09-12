prepare-test-wd:
	mkdir -p build
	cp textures.txt build/textures.txt

install:
	$(prepare-test-wd)
	GOOS=darwin go build -o build/lunarlauncher

run:
	$(prepare-test-wd)
	cd build && GOOS=darwin go run ../.
