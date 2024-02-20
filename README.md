# Docker-Backup 🚀
Docker Backup is a command-line tool for backing up Docker container volumes using BorgBackup. It simplifies the backup process, making it easy to manage and automate data backups for Docker containers.
## ⚠️ Warning
This tool is currently in its early development phase and is not recommended for use in a production environment.

## ✨ Features
- Effortless Backup: Easily backup all volumes of a specific Docker container using BorgBackup.
- Seamless Integration: Target container gets gracefully stopped during backup and restarted afterward.
- Local/Remote Backup: Backup to a Borg repository running locally or on a server, ensuring your data is stored securely.
## 🛠 Installation
### Clone Repository
```bash
git clone https://github.com/your-username/docker-backup.git
cd docker-backup
```
### Build Docker Worker Image
```bash
docker build -t worker -f build/docker/worker.Dockerfile .
```
### Configuration Adjustment
- set the container you want to backup in ```cmd/docker-backup/main.go```
- configure the repos you want to backup to in ```docker-backup.yml```
###  Build & Run Project
```bash
go run main.go
```
## 📋 Requirements
- Docker
- GoLang