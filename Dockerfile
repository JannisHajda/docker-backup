FROM ubuntu

# Update and install necessary tools
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
RUN echo "deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker CLI
RUN apt-get update && apt-get install -y docker-ce-cli

# Download the Borg binary for glibc 2.36
RUN wget -O /usr/local/bin/borg https://github.com/borgbackup/borg/releases/download/1.4.0/borg-linux-glibc236

# Set permissions for Borg binary
RUN chown root:root /usr/local/bin/borg && chmod 755 /usr/local/bin/borg

# Optional: Create a symlink for borgfs
RUN ln -s /usr/local/bin/borg /usr/local/bin/borgfs

# Install rclone
RUN curl https://rclone.org/install.sh | bash

CMD ["/bin/bash"]
