## Setup

### 1. Create binary

```bash
go build
```

### 2. start tracker

```bash
sudo ./filechangestracker
```

### 3. Check heath of of workers

`curl -s -X GET http://localhost:9000/v1/health`

### 4. Get logs

`curl -s -X GET http://localhost:9000/v1/logs`

### 5. Add new command to queue

```bash
curl -s -X POST http://localhost:9000/v1/commands \
-H "Content-Type: application/json" \
-d '{"commands":["touch /Users/user/Downloads/test/test.txt"]}'
```

NOTE:
osqueryd should be running

```bash
sudo osqueryd --socket=/Users/user/.osquery/shell.em --verbose --disable_events=false --disable_audit=false --disable_endpointsecurity=false --disable_endpointsecurity_fim=false --enable_file_events=true
```

osquery.conf

```json
"file_paths": {
    "downloads": [
      "/Users/user/Downloads/test/%%"
    ]
  }
```

Create a .pkg Installer

```bash
pkgbuild --root /usr/local/bin \
         --identifier com.test.filechangestracker \
         --version 1.0 \
         --install-location / \
         filechangestracker-1.0.pkg
```
