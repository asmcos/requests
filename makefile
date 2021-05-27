test:
	go test -v ./examples

test_local:
	go test -run TestResponseHeader -v ./examples

pkg:
	newversion.py version
	jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS


