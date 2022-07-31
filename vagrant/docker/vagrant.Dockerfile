FROM rockylinux:latest

ENV container docker



# Set repos to nju

RUN sed -e 's|^mirrorlist=|#mirrorlist=|g' \

        -e 's|^#baseurl=http://dl.rockylinux.org/$contentdir|baseurl=https://mirrors.nju.edu.cn/rocky|g' \

        -i.bak \

        etc/yum.repos.d/Rocky-*.repo



# Add epel releaes repo

RUN dnf -y install epel-release



# Perform a package update

RUN dnf -y update



# Add some familiar utilities

RUN dnf -y install procps htop grep findutils iputils iproute



# Add sshd server so we can 'vagrant ssh' later

RUN dnf -y install openssh-server openssh-clients passwd sudo;

RUN mkdir var/run/sshd

RUN ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key -N ''

RUN useradd --create-home -s /bin/bash vagrant

RUN echo -e "vagrant\nvagrant" | (passwd --stdin vagrant)

RUN echo 'vagrant ALL=(ALL) NOPASSWD: ALL' > /etc/sudoers.d/vagrant

RUN chmod 440 /etc/sudoers.d/vagrant

RUN mkdir -p /home/vagrant/.ssh

RUN chmod 700 /home/vagrant/.ssh

ADD https://raw.githubusercontent.com/hashicorp/vagrant/master/keys/vagrant.pub /home/vagrant/.ssh/authorized_keys

RUN chmod 600 /home/vagrant/.ssh/authorized_keys

RUN chown -R vagrant:vagrant /home/vagrant/.ssh

# Allow public key authentication for 'vagrant ssh' in Fedora 35

RUN sed -i 's/^#PubkeyAuthentication yes/PubkeyAuthentication yes/i' /etc/ssh/sshd_config

# This softens a crypto policy that prevents vagrant completing ssh setup

RUN sed -i 's/^Include \/etc\/crypto-policies\/back-ends\/opensshserver.config/#Include \/etc\/crypto-policies\/back-ends\/opensshserver.config/i' /etc/ssh/ssh_config.d/05-redhat.conf

# As the container isn't normally running systemd, /run/nologin needs to be removed to allow SSH

RUN rm -rf /run/nologin

# Let's install and enable nginx for fun - just to prove this works!

RUN dnf -y install nginx

RUN systemctl enable nginx

# Install the replacement systemctl command

RUN yum -y install python3

COPY src/docker-systemctl-replacement/files/docker/systemctl3.py /usr/bin/systemctl

RUN chmod 755 /usr/bin/systemctl

# COPY src/docker-systemctl-replacement/files/docker/journalctl3.py /usr/bin/journalctl

# RUN chmod 755 /usr/bin/journalctl



CMD /usr/bin/systemctl