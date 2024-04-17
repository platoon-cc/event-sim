# Platoon CLI

Command line tooling for interacting with the Platoon backend, querying of events and for debugging
client-side integrations.

## To Install

```bash
go install github.com/platoon-cc/platoon-cli@latest
```

## General Use

Platoon CLI has built in help so just running with no arguments (or with -h, --help) will display command help.

## Client Integration

When adding a client integration to your game, it can be useful to test locally before beginning to send events to the
backend. For this purpose, Platoon CLI provides a local server which can be run like this:

```bash
platoon-cli server
```

This will start listening on localhost at port 9998 (the port can be changed by passing `--port <int>` as an argument)

Then, set your client into Debug Mode and set the Debug URL to be `http://localhost:9998` (the default) and then
any events sent from your integration will be ingested into a local sqlite database.

## Querying

There are a few built-in queries. Run `platoon-cli query` to see a list - and make sure to pass the `--local` flag to run against locally gathered events.

## Building/Releasing

* Update the version

```bash
go run tools/update_version.go
git commit ...
git push ...
./tools/push_tag.sh
```
