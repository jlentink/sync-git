FROM alpine/git:latest
ADD sync-git /usr/bin/sync-git
WORKDIR /repos
ENTRYPOINT ["/usr/bin/sync-git", "sync", "-p"]