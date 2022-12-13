FROM alpine/git:latest
ADD sync-git /usr/bin/sync-git
ENTRYPOINT ["/usr/bin/sync-git", "sync", "-p"]