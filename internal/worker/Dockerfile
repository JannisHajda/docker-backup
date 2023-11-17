FROM debian:stretch-slim

# Install dependencies (borgbackup)
RUN apt-get update && apt-get install -y \
    borgbackup 

# RUN command (borg backup)
RUN borg create --stats --progress --compression lz4 --exclude-caches --exclude-if-present .nobackup /output::'{now:%Y-%m-%d_%H:%M:%S}' /input 