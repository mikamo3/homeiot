FROM mcr.microsoft.com/devcontainers/go:0-1-bullseye
ARG USERNAME=vscode
ARG USER_UID=1000
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | \
    apt-key add - && echo "deb https://dl.yarnpkg.com/debian/ stable main" | \
    tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install -y \
    bluez \
    dbus \
    sudo \
    ncat \
    mosquitto \
    mosquitto-clients
RUN  echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
RUN sudo usermod -a -G bluetooth ${USERNAME}
