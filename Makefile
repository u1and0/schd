install:
	GOOS=windows GOARCH=amd64 go build -o release/schd.exe
	cd static && npx tsc || cd ..
	rsync -auv static release
	rsync -auv template release
