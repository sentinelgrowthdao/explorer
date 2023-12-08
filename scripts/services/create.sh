#!/bin/bash

app_directory="./cmd"
systemd_dir="/etc/systemd/system"
app_template="[Unit]
Description=Sentinel Explorer APP_NAME Daemon
After=network-online.target

[Service]
User=root
ExecStart=APP_NAME
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target"

for app_name in "$app_directory"/*; do
  app_name=$(basename "$app_name")
  service_file="$systemd_dir/${app_name}.service"
  modified_template="${app_template//APP_NAME/$app_name}"

  ln -s "${GOPATH}/bin/${app_name}" "/usr/local/bin/${app_name}"
  if echo "$modified_template" > "$service_file"; then
    echo "Created $service_file"
  fi
done
