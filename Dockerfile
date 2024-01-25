FROM ubuntu
RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    borgbackup

CMD ["tail", "-f", "/dev/null"]
