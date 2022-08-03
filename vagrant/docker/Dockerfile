FROM centos:centos7
MAINTAINER liz <rick.lizhan@gmail.com>

RUN yum -y update &&\
    yum -y install kde-l10n-Chinese telnet &&\
    yum -y reinstall glibc-common &&\
    yum -y install openssh-server openssh-clients passwd sudo &&\
    yum clean all &&\
    rm -rf /tmp/* rm -rf /var/cache/yum/* &&\
    localedef -i en_US -f UTF-8 en_US.UTF-8 &&\
    localedef -c -f UTF-8 -i zh_CN zh_CN.utf8

RUN mkdir /var/run/sshd
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

ENTRYPOINT ["/usr/sbin/sshd", "-D"]
