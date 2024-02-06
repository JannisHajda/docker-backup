# docker-backup

Docker Backup is a command-line tool for backing up Docker container volumes using BorgBackup. It simplifies the backup process, making it easy to manage and automate data backups for Docker containers.

## Features 
- backup all volumes of a specific docker container
  using borgbackup to a volume (target container gets stopped during and restarted after backup)
```bash
docker-backup backup --container <container> --passphrase <passphrase> --output <output-volume>
```

## Installation
- Clone this repository
- Build the docker worker image
  ```bash
  docker build -t worker .
  ```
- Build the go binary
  ```bash
  go build -o docker-backup cmd/docker-backup/main.go
  ```

## Requirements
- Docker
- GoLang
