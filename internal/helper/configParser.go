package helper

import (
	"docker-backup/interfaces"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Repos map[string]BackupConfig `yaml:"repos"`
}

type BackupConfig struct {
	Type       string `yaml:"type"`
	VolumeName string `yaml:"volumeName"`
	Path       string `yaml:"path"`
	Passphrase string `yaml:"borg-passphrase"`
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	SSHKey     string `yaml:"ssh-key"`
}

func ParseConfigFile(filePath string) ([]interfaces.LocalBackup, []interfaces.RemoteBackup, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, nil, err
	}

	var localBackups []interfaces.LocalBackup
	var remoteBackups []interfaces.RemoteBackup

	for _, backupConfig := range config.Repos {
		switch backupConfig.Type {
		case "local":
			localBackups = append(localBackups, interfaces.LocalBackup{
				Backup: interfaces.Backup{
					Path:       backupConfig.Path,
					Passphrase: backupConfig.Passphrase,
					Keyfile:    "",
				},
				VolumeName: backupConfig.VolumeName,
			})
		case "remote":
			remoteBackups = append(remoteBackups, interfaces.RemoteBackup{
				Backup: interfaces.Backup{
					Path:       backupConfig.Path,
					Passphrase: backupConfig.Passphrase,
					Keyfile:    "",
				},
				User:   backupConfig.User,
				Host:   backupConfig.Host,
				SSHKey: backupConfig.SSHKey,
			})
		default:
			panic("Unknown backup type: " + backupConfig.Type)
		}
	}

	return localBackups, remoteBackups, nil
}
