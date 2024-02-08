FROM ubuntu
RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    borgbackup openssh-client iputils-ping

ENV LANG en_US.UTF-8

# mark all subfolders in /input as volumes
VOLUME /source
VOLUME /backups

CMD ["tail", "-f", "/dev/null"]
