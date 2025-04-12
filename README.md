# stress-test

```bash
docker build -t stress-test .
docker run stress-test --url=http://google.com --requests=50 --concurrency=10
```

Small project that simulates a stress test, using Go.

This has been coded during a post-graduation that I've done in Go, and its not really meant to be a project to be used daily.

The code is here for future reference from me or anyone who needs it.
