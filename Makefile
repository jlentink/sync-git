.PHONY: clean release test-release shell lint

clean:
	@-rm -rf dist cmd/sync-git/sync-git

release: clean
	goreleaser

test-release:
	goreleaser --skip-publish --skip-validate --rm-dist --auto-snapshot

shell:
	docker run -ti --entrypoint /bin/sh jlentink/sync-git:latest

lint:
	golangci-lint run
