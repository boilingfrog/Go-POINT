# Modified by Matthew Warman <mcwarman@gmail.com> 
# Based on project with in CenOs/CentOS-Dockerfiles
# https://github.com/CentOS/CentOS-Dockerfiles/blob/master/ssh/centos7
#
# Which was "ported" by Adam Miller <maxamillion@fedoraproject.org> from
#   https://github.com/fedora-cloud/Fedora-Dockerfiles
#
# Originally written for Fedora-Dockerfiles by
#   scollier <scollier@redhat.com>

FROM centos:centos7
MAINTAINER liz <rick.lizham@gmail.com>
RUN yum -y update; yum clean all
RUN yum -y install openssh-server openssh-clients passwd sudo; yum clean all
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
