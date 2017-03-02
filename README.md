# Remote File Resource

Downloads a file from a remote server.

## Source Configuration

* `uri`: *Required.* The location of the file to download. This URI must return an ETag to use for versioning.

## Behavior

### `check`: Does a HEAD request to the provided URI to see if the ETag has changed.
