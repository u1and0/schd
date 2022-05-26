install:
	GOOS=windows GOARCH=amd64 go build -o release/schd.exe
	cd static && npx tsc || cd ..
	rsync -auv --delete static release
	rsync -auv --delete template release
