test:
	go test -v requests_test.go requests.go

pkg:
	newversion.py version
	jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS

