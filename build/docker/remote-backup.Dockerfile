FROM ubuntu
RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    borgbackup openssh-client openssh-server iputils-ping

ENV LANG en_US.UTF-8

# mark all subfolders in /input as volumes
VOLUME /source
VOLUME /backups

RUN useradd -m -s /bin/bash borg
RUN echo "borg:passwd" | chpasswd

RUN chown -R borg:borg /home/borg

COPY remote-backup-entrypoint.sh /usr/local/bin/remote-backup-entrypoint.sh
RUN chmod +x /usr/local/bin/remote-backup-entrypoint.sh

CMD ["/usr/local/bin/remote-backup-entrypoint.sh"]