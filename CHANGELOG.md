# CHANGELOG

## Release builds

Releases are build with github actions.

Binaries are provided to download for linux for platforms 
- amd64
- arm64
- arm

Docker container are build for the same platforms:

```
ghcr.io/ricoschulte/sysclienttester:latest
ghcr.io/ricoschulte/sysclienttester:<tag>
ghcr.io/ricoschulte/sysclienttester:v0.0.2
```

## Default idientities included

The default identity file is avaiable inside the container build at /identities.json