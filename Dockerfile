FROM ubuntu
RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    borgbackup

ENV LANG en_US.UTF-8

CMD ["tail", "-f", "/dev/null"]
