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
Copy code
docker build -t worker .
```
### Configuration Adjustment
Navigate to main.go and adapt the configuration based on your requirements.
###Build & Run Project
```bash
go run main.go
```
## 📋 Requirements
- Docker
- GoLang