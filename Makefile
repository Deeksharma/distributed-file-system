# target: prerequisites (space separated, needs to be present before the target)
build:
	go build -o bin/fs # these commands should be tab separated

run: build
	@./bin/fs

test: build run
	go test ./... -v


# send a file
# hash the file
# key
#  add interface func to transform the key
# subfolders in pars of 2 ans tore the data somewhere there