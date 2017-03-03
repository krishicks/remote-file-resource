# Remote File Resource

Downloads a file from a remote server.

## Source Configuration

* `uri`: *Required.* The location of the file to download. This URI must return an ETag to use for versioning.

## Behavior

### `check`: Does a HEAD request to the provided URI to see if the ETag has changed.

### `in`: Download the item at the configured URI.

Places the downloaded file in the destination with the provided filename.

#### Parameters

* `filename`: *Required.* The filename that will be used for the downloaded file.

### `out`: Not implemented.
