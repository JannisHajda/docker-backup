FROM ubuntu
RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    borgbackup openssh-client sshpass iputils-ping

ENV LANG en_US.UTF-8

# mark all subfolders in /input as volumes
VOLUME /source
VOLUME /backups

RUN mkdir -p /root/.ssh
RUN ssh-keygen -t ed25519 -f /root/.ssh/id_ed25519 -C "worker" -N ""

COPY worker-entrypoint.sh /usr/local/bin/worker-entrypoint.sh
RUN chmod +x /usr/local/bin/worker-entrypoint.sh

CMD ["tail", "-f", "/dev/null"]
