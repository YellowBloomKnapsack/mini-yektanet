#!/usr/bin/env bash

# Declare variables
SERVER_USER=
SERVER_IP=
SERVER_PORT=
REMOTE_DIR="/home/$SERVER_USER/mini-yektanet"

# Function to perform SCP
perform_scp() {
    local source_file="$1"
    local dest_path="$2"
    scp -P $SERVER_PORT "$source_file" "$SERVER_USER@$SERVER_IP:$dest_path"
}

# Perform SCP for each .env file
perform_scp "panel/.env" "$REMOTE_DIR/panel/.env"
perform_scp "adserver/.env" "$REMOTE_DIR/adserver/.env"
perform_scp "common/.env" "$REMOTE_DIR/common/.env"
perform_scp "eventserver/.env" "$REMOTE_DIR/eventserver/.env"
perform_scp "publisherwebsite/.env" "$REMOTE_DIR/publisherwebsite/.env"
