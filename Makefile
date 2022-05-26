install:
	GOOS=windows GOARCH=amd64 go build -o release/schd.exe
	go run main.go -v > release/VERSION.txt
	cd static && npx tsc || cd ..
	rsync -auv --delete static release
	rsync -auv --delete template release
