## Setup

### 1. Config

- Update the config.yaml to setup tracking directory and osquery socket_path
- Find your socket path
  - run in terminal: `osqueryi`
  - run query `select value from osquery_flags where name = 'extensions_socket';`
  - value returned is your socket path
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

### 4. Get logs

`curl -s -X GET http://localhost:9000/v1/logs`

### 5. Add new command to queue

```bash
curl -s -X POST http://localhost:9000/v1/commands \
-H "Content-Type: application/json" \
-d '{"commands":["touch /Users/{USERNAME}/Downloads/test/test.txt"]}'
```

NOTE:
osqueryd should be running

- Update osquery.conf

````json
{
  "file_paths": {
    "downloads": ["/Users/{USERNAME}/Downloads/test/%%"]
  }
}

- macOS:

```bash
sudo osqueryd --socket=/Users/{USERNAME}/.osquery/shell.em --verbose --disable_events=false --disable_audit=false --disable_endpointsecurity=false --disable_endpointsecurity_fim=false --enable_file_events=true
````

````

- Windows: (_as administrator_)

```bash
osqueryd --socket=C:\Users\{USERNAME}\.osquery\shell.em --verbose --disable_events=false --enable_ntfs_event_publisher=true --enable_powershell_events_subscriber=true --enable_windows_events_publisher=true --enable_windows_events_subscriber=true

````
