#!/bin/bash
set -e

# Check if osquery is installed using brew
if ! brew list --cask | grep -q "osquery" && [ ! -f "/opt/osquery/lib/osquery.app/Contents/MacOS/osqueryd" ]; then
    echo "osquery not found and osqueryd binary is missing. Installing osquery via Homebrew..."
    brew install --cask osquery
    echo "osquery installed successfully."
else
    echo "osquery is already installed or the osqueryd binary exists."
fi

JSON_CONTENT="{\"file_paths\": {\"downloads\": [\"$HOME/Downloads/%%\"]}}"

echo $JSON_CONTENT | sudo tee /var/osquery/osquery.conf > /dev/null

echo "Configuration added to /var/osquery/osquery.conf"


# YAML content to be written to config.yaml
YAML_CONTENT=$(cat <<EOF
directory: '$HOME/Downloads/'
check_frequency: 1 # in seconds
reporting_api: 'http://api.external.com/v1/report'
http_port: '9000'
socket_path: '/var/osquery/osquery.em'
mongo_uri: 'mongodb://user:password@localhost:27017'
EOF
)


CONFIG_FILE="config.yaml"

echo "$YAML_CONTENT" | sudo tee "$CONFIG_FILE" > /dev/null

echo "YAML configuration added to $CONFIG_FILE"

echo "setup complete!!"
