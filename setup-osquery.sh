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

# Add /opt/osquery/lib/osquery.app/Contents/MacOS to PATH in the user's shell profile
SHELL_PROFILE=""

if [ "$SHELL" = "/bin/bash" ]; then
    SHELL_PROFILE="$HOME/.bash_profile"
elif [ "$SHELL" = "/bin/zsh" ]; then
    SHELL_PROFILE="$HOME/.zprofile"
else
    echo "Unsupported shell. Please add osqueryd path manually to your shell profile."
fi


# Add to PATH if not already present
if ! grep -q "/opt/osquery/lib/osquery.app/Contents/MacOS" "$SHELL_PROFILE"; then
    echo -e '\n\n# osquery path\nexport PATH=$PATH:/opt/osquery/lib/osquery.app/Contents/MacOS' | sudo tee -a "$SHELL_PROFILE" > /dev/null
    echo "PATH updated in $SHELL_PROFILE"
else
    echo "PATH already contains /opt/osquery/lib/osquery.app/Contents/MacOS"
fi

if [ "$SHELL" = "/bin/bash" ]; then
    source $SHELL_PROFILE
elif [ "$SHELL" = "/bin/zsh" ]; then
    /bin/zsh -c "source $SHELL_PROFILE"
else
    echo "Unsupported shell. Please add osqueryd path manually to your shell profile and run `source $SHELL_PROFILE`"
fi


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