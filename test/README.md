# dev-infra
```
docker run -it -v $GOPATH/src/github.com/paulhindemith/dev-infra:/go/src/github.com/paulhindemith/dev-infra -v $GOPATH/pkg/dep:/go/pkg/dep paulhindemith/golang-dev:1.13

cd dev-infra
./test/presubmit-tests.sh

```
