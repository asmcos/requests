msg ?= ''
test:
	go test -v ./examples

test_local:
	go test -run TestResponseHeader -v ./examples

.ONESHELL:
t:
	if [[ "$(msg)" = "" ]] ; then echo "Usage: make pkg msg='commit msg'";exit 20; fi

gitcheck:
	if [[ "$(msg)" = "" ]] ; then echo "Usage: make pkg msg='commit msg'";exit 20; fi

.ONESHELL:
pkg: gitcheck test
	newversion.py version
	git commit -am "$(msg)"
	jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag "$$v" && git push origin "$$v"
