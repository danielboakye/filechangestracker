## Setup

### 1. Config

- setup config and osquery
  on mac

```bash
make setup/osquery/mac
```

- setup mongo db in docker

```bash
make logsdb
```

### 2. Build and start

```bash
make start
```

### 3. Check heath of of workers

`curl -s -X GET http://localhost:9000/v1/health`

### 4. Add new command to queue

```bash
touch $HOME/Downloads/test1.txt
```

```bash
curl -s -X POST http://localhost:9000/v1/commands \
-H "Content-Type: application/json" \
-d "{\"commands\":[\"touch $HOME/Downloads/test2.txt\"]}"
```

### 5. Get logs

- wait 5 seconds and run command below

`curl -s -X GET http://localhost:9000/v1/logs\?limit=2`

NOTES

osqueryd flags

- macOS:

  ```bash
  sudo osqueryd --verbose --disable_events=false --disable_audit=false --disable_endpointsecurity=false --disable_endpointsecurity_fim=false --enable_file_events=true
  ```

- Windows: (_as administrator_)

  ```bash
  osqueryd --verbose --disable_events=false --enable_ntfs_event_publisher=true --enable_powershell_events_subscriber=true --enable_windows_events_publisher=true --enable_windows_events_subscriber=true

  ```
