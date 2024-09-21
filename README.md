## Setup

### 1. Config

Updated config.yaml to setup tracking directory and osquery socket_path

### 2. Build and start

```bash
go build && sudo ./filechangestracker
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

- Windows: (_as administrator_)

```bash
osqueryd --socket=C:\Users\user\.osquery\shell.em --verbose --disable_events=false --enable_ntfs_event_publisher=true --enable_powershell_events_subscriber=true --enable_windows_events_publisher=true --enable_windows_events_subscriber=true

```

- macOS:

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
