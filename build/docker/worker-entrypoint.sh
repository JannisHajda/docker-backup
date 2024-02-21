#!/bin/bash

# Add KnownHosts
ssh-keyscan -H remote-backup >> /root/.ssh/known_hosts

# Copy SSH key to remote-backup
sshpass -p passwd ssh-copy-id -i /root/.ssh/id_ed25519 borg@remote-backup

# Continue with the original CMD
tail -f /dev/null
