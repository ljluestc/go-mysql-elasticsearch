# Reindex API Support

This document describes the Elasticsearch Reindex API support in `go-mysql-elasticsearch`.

## Overview

We expose the Elasticsearch `_reindex` API via a new HTTP endpoint on the status server. This allows for triggering manual reindexing operations directly.

## Usage

To trigger a reindex, send a `POST` request to the `/reindex` endpoint of the status server (default port `12800`). The body of the request should contain the standard Elasticsearch reindex JSON payload.

### Example

```bash
curl -X POST http://127.0.0.1:12800/reindex -d '{
  "source": {
    "index": "old_index"
  },
  "dest": {
    "index": "new_index"
  }
}'
```

## Implementation Details

The implementation involves:
- **Elasticsearch Client**: A `Reindex` method has been added to the internal client.
- **HTTP Handler**: A new handler `/reindex` accepts JSON requests and proxies them to the Elasticsearch `_reindex` API.
