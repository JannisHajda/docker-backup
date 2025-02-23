FROM ubuntu:22.04

# Switch to bash with pipefail so the build fails if install.sh fails
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release \
    wget

# Add Docker’s official GPG key
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -

# Set up the stable Docker repository
RUN echo "deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list

# Install Docker CLI
RUN apt-get update && apt-get install -y docker-ce-cli

# Download and set up Borg
RUN wget -O /usr/local/bin/borg \
      https://github.com/borgbackup/borg/releases/download/1.4.0/borg-linux-glibc236 && \
    chown root:root /usr/local/bin/borg && \
    chmod 755 /usr/local/bin/borg

# Optional symlink for borgfs
RUN ln -s /usr/local/bin/borg /usr/local/bin/borgfs

# Install rclone. If the script fails, the build fails due to pipefail
RUN curl https://rclone.org/install.sh | bash && \
    which rclone && \
    rclone version

CMD ["/bin/bash"]
