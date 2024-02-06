# docker-backup

Docker Backup is a command-line tool for backing up Docker container volumes using BorgBackup. It simplifies the backup process, making it easy to manage and automate data backups for Docker containers.

## Features (so far)
- backup all volumes of a specific docker container
  using borgbackup to a volume

## Installation
- Clone this repository
- Build the docker worker image
  ```bash
  docker build -t worker .
  ```
- (currently using `main.go`as entrypoint, cli will follow soon)