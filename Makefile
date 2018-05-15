snapshot: # create a snapshot
	goreleaser --snapshot --rm-dist --skip-validate --skip-publish

release: # create a release build
	goreleaser --rm-dist

install: # just compile and install
	go install

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help