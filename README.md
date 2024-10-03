# SFS (Stupid File Storage)


# SFSP (Stupid File Storage Protocol)
> v0.1.0

## Send chunk 

### Request
Format:

```
*<filename_size><filename><id><size><data>
```

Where:
- `filename_size` is a little-endian uint64
- `id` is a little-endian uint64
- `size` is a little-endian uint64
- `filename` is []byte with len of `filename_size`, containing the name
  of uploading file
- `data` is []byte with len of `size` containing the `id`'th chunk of file

### Responce

```
<code><msg_size>[<msg>]
```

Where:
- `code` is the little-endian uint64 status code (see [Status Codes](#status-codes))
- `msg_size` is the little-endian uint64 representing the size of following `msg`. If `msg` not present, the `msg_size` will be `0`

## Receive chunk

### Request
Request specified chunk of the file.

Format:

```
/<filename_size><filename><id>
```

Where:
- `filename_size` is a little-endian uint64
- `id` is a little-endian uint64 ID of file chunk
- `filename` is []byte with len of `filename_size`, containing the name
  of uploading file

### Responce

Depends on the `code`

#### `code` is `OK`:

```
<code><filename_size><filename><id><size><data>
```

The part after `code` is identical to [send chunk request](#send-chunk)

#### `code` is `NOT_FOUND`:

```
<code>
```


#### `code` is `INTERNAL`:

```
<code><msg_size>[<msg>]
```

Where:
- `msg_size` is the little-endian uint64 representing the size of following `msg`. If `msg` not present, the `msg_size` will be `0`

## Receive all chunks ids stored in node

### Request
Request specific chunks of the file.

Format:

```
%<filename_size><filename>
```

Where:
- `filename_size` is a little-endian uint64
- `filename` is []byte with len of `filename_size`, containing the name
  of uploading file
  
### Response

#### `code` is `OK`:

```
<code><count>[<...ids>]
```

Where:
- `count` is a little-endian uint64 representing count of following chunk IDs
- `...ids` is a sequence (len = `count`) of little-endian uint64, representing ids of chunks stored in node

If there is no chunks in the node, the `count` will be `0` and there will be no ids after that. The `code` still will be `OK`

#### `code` is `INTERNAL`:

```
<code><msg_size>[<msg>]
```

Where:
- `msg_size` is the little-endian uint64 representing the size of following `msg`. If `msg` not present, the `msg_size` will be `0`

----------------------------------

## Invalid Request
If the request will not match to any of specified requests, the server will return `INVALID_REQ` code in following format:

```
<code><msg_size>[<msg>]
```

Where:
- `msg_size` is the little-endian uint64 representing the size of following `msg`. If `msg` not present, the `msg_size` will be `0`

## Status Codes
Status code is a little-endian uint64 with following meanings:
- `10` - OK
- `20` - NOT_FOUND
- `21` - INVALID_REQ
- `30` - INTERNAL

## FAQ
### Why are you using little-endian uint64 for everything?
i don't know
