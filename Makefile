# make cyclic
cyclic:
	go install github.com/elza2/go-cyclic@latest
	go-cyclic run --dir . filter *_test.go
