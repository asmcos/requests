msg?=''
test:
	go test -v ./examples

test_local:
	go test -run TestResponseHeader -v ./examples

.ONESHELL:
t:
	v=`cat version` && echoraw $$v
	v=`cat version` && git tag $$v && git push origin $$v

pkg:
	echoraw "$$msg"
	exit
	newversion.py version
	git commit -am "$$msg"
	jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag $$v && git push origin $$v



