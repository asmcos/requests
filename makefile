msg ?= ''
test:
	go test -v ./examples

.ONESHELL:
gitcheck:
	if [[ "$(msg)" = "" ]] ; then echo "Usage: make pkg msg='commit msg'";exit 20; fi


.ONESHELL:
pkg: gitcheck test
	hash  newversion.py && newversion.py version
	git commit -am "$(msg)"
	jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag "$$v" && git push origin "$$v"
