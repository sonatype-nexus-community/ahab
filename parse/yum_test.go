// Copyright 2019 Sonatype Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package parse

import (
	"strings"
	"testing"
)

var dnfList = `Installed Packages
acl.x86_64                                                                  2.2.53-1.el8                                                             @System
audit-libs.x86_64                                                           3.0-0.17.20191104git1c2f876.el8                                          @System
basesystem.noarch                                                           11-5.el8                                                                 @System
bash.x86_64                                                                 4.4.19-10.el8                                                            @System
bind-export-libs.x86_64                                                     32:9.11.13-5.el8_2                                                       @System
binutils.x86_64                                                             2.30-73.el8                                                              @System
bzip2-libs.x86_64                                                           1.0.6-26.el8                                                             @System
ca-certificates.noarch                                                      2020.2.41-80.0.el8_2                                                     @BaseOS
centos-gpg-keys.noarch                                                      8.2-2.2004.0.1.el8                                                       @System
centos-release.x86_64                                                       8.2-2.2004.0.1.el8                                                       @System
centos-repos.x86_64                                                         8.2-2.2004.0.1.el8                                                       @System
chkconfig.x86_64                                                            1.11-1.el8                                                               @System
coreutils-single.x86_64                                                     8.30-7.el8_2.1                                                           @System
cpio.x86_64                                                                 2.12-8.el8                                                               @System
cracklib.x86_64                                                             2.9.6-15.el8                                                             @System
crypto-policies.noarch                                                      20191128-2.git23e1bf1.el8                                                @System
cryptsetup-libs.x86_64                                                      2.2.2-1.el8                                                              @System
curl.x86_64                                                                 7.61.1-12.el8                                                            @System
cyrus-sasl-lib.x86_64                                                       2.1.27-1.el8                                                             @System
dbus.x86_64                                                                 1:1.12.8-10.el8_2                                                        @System
dbus-common.noarch                                                          1:1.12.8-10.el8_2                                                        @System
dbus-daemon.x86_64                                                          1:1.12.8-10.el8_2                                                        @System
dbus-libs.x86_64                                                            1:1.12.8-10.el8_2                                                        @System
dbus-tools.x86_64                                                           1:1.12.8-10.el8_2                                                        @System
device-mapper.x86_64                                                        8:1.02.169-3.el8                                                         @System
device-mapper-libs.x86_64                                                   8:1.02.169-3.el8                                                         @System
dhcp-client.x86_64                                                          12:4.3.6-40.el8                                                          @System
dhcp-common.noarch                                                          12:4.3.6-40.el8                                                          @System
dhcp-libs.x86_64                                                            12:4.3.6-40.el8                                                          @System
dnf.noarch                                                                  4.2.17-7.el8_2                                                           @System
dnf-data.noarch                                                             4.2.17-7.el8_2                                                           @System
dracut.x86_64                                                               049-70.git20200228.el8                                                   @System
dracut-network.x86_64                                                       049-70.git20200228.el8                                                   @System
dracut-squash.x86_64                                                        049-70.git20200228.el8                                                   @System
elfutils-default-yama-scope.noarch                                          0.178-7.el8                                                              @System
elfutils-libelf.x86_64                                                      0.178-7.el8                                                              @System
elfutils-libs.x86_64                                                        0.178-7.el8                                                              @System
ethtool.x86_64                                                              2:5.0-2.el8                                                              @System
expat.x86_64                                                                2.2.5-3.el8                                                              @System
file-libs.x86_64                                                            5.33-13.el8                                                              @System
filesystem.x86_64                                                           3.8-2.el8                                                                @System
findutils.x86_64                                                            1:4.6.0-20.el8                                                           @System
gawk.x86_64                                                                 4.2.1-1.el8                                                              @System
gdbm.x86_64                                                                 1:1.18-1.el8                                                             @System
gdbm-libs.x86_64                                                            1:1.18-1.el8                                                             @System
glib2.x86_64                                                                2.56.4-8.el8                                                             @System
glibc.x86_64                                                                2.28-101.el8                                                             @System
glibc-common.x86_64                                                         2.28-101.el8                                                             @System
glibc-minimal-langpack.x86_64                                               2.28-101.el8                                                             @System
gmp.x86_64                                                                  1:6.1.2-10.el8                                                           @System
gnupg2.x86_64                                                               2.2.9-1.el8                                                              @System
gnutls.x86_64                                                               3.6.8-11.el8_2                                                           @System
gpgme.x86_64                                                                1.10.0-6.el8.0.1                                                         @System
grep.x86_64                                                                 3.1-6.el8                                                                @System
gzip.x86_64                                                                 1.9-9.el8                                                                @System
hostname.x86_64                                                             3.20-6.el8                                                               @System
ima-evm-utils.x86_64                                                        1.1-5.el8                                                                @System
info.x86_64                                                                 6.5-6.el8                                                                @System
ipcalc.x86_64                                                               0.2.4-4.el8                                                              @System
iproute.x86_64                                                              5.3.0-1.el8                                                              @System
iptables-libs.x86_64                                                        1.8.4-10.el8_2.1                                                         @System
iputils.x86_64                                                              20180629-2.el8                                                           @System
json-c.x86_64                                                               0.13.1-0.2.el8                                                           @System
kexec-tools.x86_64                                                          2.0.20-14.el8                                                            @System
keyutils-libs.x86_64                                                        1.5.10-6.el8                                                             @System
kmod.x86_64                                                                 25-16.el8                                                                @System
kmod-libs.x86_64                                                            25-16.el8                                                                @System
krb5-libs.x86_64                                                            1.17-18.el8                                                              @System
langpacks-en.noarch                                                         1.0-12.el8                                                               @System
less.x86_64                                                                 530-1.el8                                                                @System
libacl.x86_64                                                               2.2.53-1.el8                                                             @System
libarchive.x86_64                                                           3.3.2-8.el8_1                                                            @System
libassuan.x86_64                                                            2.5.1-3.el8                                                              @System
libattr.x86_64                                                              2.4.48-3.el8                                                             @System
libblkid.x86_64                                                             2.32.1-22.el8                                                            @System
libcap.x86_64                                                               2.26-3.el8                                                               @System
libcap-ng.x86_64                                                            0.7.9-5.el8                                                              @System
libcom_err.x86_64                                                           1.45.4-3.el8                                                             @System
libcomps.x86_64                                                             0.1.11-4.el8                                                             @System
libcurl-minimal.x86_64                                                      7.61.1-12.el8                                                            @System
libdb.x86_64                                                                5.3.28-37.el8                                                            @System
libdb-utils.x86_64                                                          5.3.28-37.el8                                                            @System
libdnf.x86_64                                                               0.39.1-6.el8_2                                                           @System
libfdisk.x86_64                                                             2.32.1-22.el8                                                            @System
libffi.x86_64                                                               3.1-21.el8                                                               @System
libgcc.x86_64                                                               8.3.1-5.el8.0.2                                                          @System
libgcrypt.x86_64                                                            1.8.3-4.el8                                                              @System
libgpg-error.x86_64                                                         1.31-1.el8                                                               @System
libidn2.x86_64                                                              2.2.0-1.el8                                                              @System
libkcapi.x86_64                                                             1.1.1-16_1.el8                                                           @System
libkcapi-hmaccalc.x86_64                                                    1.1.1-16_1.el8                                                           @System
libksba.x86_64                                                              1.3.5-7.el8                                                              @System
libmetalink.x86_64                                                          0.1.3-7.el8                                                              @System
libmnl.x86_64                                                               1.0.4-6.el8                                                              @System
libmodulemd1.x86_64                                                         1.8.16-0.2.8.2.1                                                         @System
libmount.x86_64                                                             2.32.1-22.el8                                                            @System
libnghttp2.x86_64                                                           1.33.0-3.el8_2.1                                                         @System
libnsl2.x86_64                                                              1.2.0-2.20180605git4a062cf.el8                                           @System
libpcap.x86_64                                                              14:1.9.0-3.el8                                                           @System
libpwquality.x86_64                                                         1.4.0-9.el8                                                              @System
librepo.x86_64                                                              1.11.0-2.el8                                                             @System
libreport-filesystem.x86_64                                                 2.9.5-10.el8                                                             @System
libseccomp.x86_64                                                           2.4.1-1.el8                                                              @System
libselinux.x86_64                                                           2.9-3.el8                                                                @System
libsemanage.x86_64                                                          2.9-2.el8                                                                @System
libsepol.x86_64                                                             2.9-1.el8                                                                @System
libsigsegv.x86_64                                                           2.11-5.el8                                                               @System
libsmartcols.x86_64                                                         2.32.1-22.el8                                                            @System
libsolv.x86_64                                                              0.7.7-1.el8                                                              @System
libstdc++.x86_64                                                            8.3.1-5.el8.0.2                                                          @System
libtasn1.x86_64                                                             4.13-3.el8                                                               @System
libtirpc.x86_64                                                             1.1.4-4.el8                                                              @System
libunistring.x86_64                                                         0.9.9-3.el8                                                              @System
libusbx.x86_64                                                              1.0.22-1.el8                                                             @System
libutempter.x86_64                                                          1.1.6-14.el8                                                             @System
libuuid.x86_64                                                              2.32.1-22.el8                                                            @System
libverto.x86_64                                                             0.3.0-5.el8                                                              @System
libxcrypt.x86_64                                                            4.1.1-4.el8                                                              @System
libxml2.x86_64                                                              2.9.7-7.el8                                                              @System
libyaml.x86_64                                                              0.1.7-5.el8                                                              @System
libzstd.x86_64                                                              1.4.2-2.el8                                                              @System
lua-libs.x86_64                                                             5.3.4-11.el8                                                             @System
lz4-libs.x86_64                                                             1.8.1.2-4.el8                                                            @System
lzo.x86_64                                                                  2.08-14.el8                                                              @System
mpfr.x86_64                                                                 3.1.6-1.el8                                                              @System
ncurses-base.noarch                                                         6.1-7.20180224.el8                                                       @System
ncurses-libs.x86_64                                                         6.1-7.20180224.el8                                                       @System
nettle.x86_64                                                               3.4.1-1.el8                                                              @System
npth.x86_64                                                                 1.5-4.el8                                                                @System
openldap.x86_64                                                             2.4.46-11.el8_1                                                          @System
openssl-libs.x86_64                                                         1:1.1.1c-15.el8                                                          @System
p11-kit.x86_64                                                              0.23.14-5.el8_0                                                          @System
p11-kit-trust.x86_64                                                        0.23.14-5.el8_0                                                          @System
pam.x86_64                                                                  1.3.1-8.el8                                                              @System
pcre.x86_64                                                                 8.42-4.el8                                                               @System
pcre2.x86_64                                                                10.32-1.el8                                                              @System
platform-python.x86_64                                                      3.6.8-23.el8                                                             @System
platform-python-setuptools.noarch                                           39.2.0-5.el8                                                             @System
popt.x86_64                                                                 1.16-14.el8                                                              @System
procps-ng.x86_64                                                            3.3.15-1.el8                                                             @System
python3-dnf.noarch                                                          4.2.17-7.el8_2                                                           @System
python3-gpg.x86_64                                                          1.10.0-6.el8.0.1                                                         @System
python3-hawkey.x86_64                                                       0.39.1-6.el8_2                                                           @System
python3-libcomps.x86_64                                                     0.1.11-4.el8                                                             @System
python3-libdnf.x86_64                                                       0.39.1-6.el8_2                                                           @System
python3-libs.x86_64                                                         3.6.8-23.el8                                                             @System
python3-pip-wheel.noarch                                                    9.0.3-16.el8                                                             @System
python3-rpm.x86_64                                                          4.14.2-37.el8                                                            @System
python3-setuptools-wheel.noarch                                             39.2.0-5.el8                                                             @System
readline.x86_64                                                             7.0-10.el8                                                               @System
rootfiles.noarch                                                            8.1-22.el8                                                               @System
rpm.x86_64                                                                  4.14.2-37.el8                                                            @System
rpm-build-libs.x86_64                                                       4.14.2-37.el8                                                            @System
rpm-libs.x86_64                                                             4.14.2-37.el8                                                            @System
sed.x86_64                                                                  4.5-1.el8                                                                @System
setup.noarch                                                                2.12.2-5.el8                                                             @System
shadow-utils.x86_64                                                         2:4.6-8.el8                                                              @System
snappy.x86_64                                                               1.1.7-5.el8                                                              @System
sqlite-libs.x86_64                                                          3.26.0-6.el8                                                             @System
squashfs-tools.x86_64                                                       4.3-19.el8                                                               @System
systemd.x86_64                                                              239-31.el8_2.2                                                           @System
systemd-libs.x86_64                                                         239-31.el8_2.2                                                           @System
systemd-pam.x86_64                                                          239-31.el8_2.2                                                           @System
systemd-udev.x86_64                                                         239-31.el8_2.2                                                           @System
tar.x86_64                                                                  2:1.30-4.el8                                                             @System
tzdata.noarch                                                               2020a-1.el8                                                              @System
util-linux.x86_64                                                           2.32.1-22.el8                                                            @System
vim-minimal.x86_64                                                          2:8.0.1763-13.el8                                                        @System
which.x86_64                                                                2.21-12.el8                                                              @BaseOS
xz.x86_64                                                                   5.2.4-3.el8                                                              @System
xz-libs.x86_64                                                              5.2.4-3.el8                                                              @System
yum.noarch                                                                  4.2.17-7.el8_2                                                           @System
zlib.x86_64                                                                 1.2.11-13.el8                                                            @System`

// generate CLI package list via:
// # yum list installed
var yumList = `Loaded plugins: fastestmirror, ovl
Installed Packages
MAKEDEV.x86_64                   3.24-6.el6                        @CentOS/6.10
audit-libs.x86_64                2.4.5-6.el6                       @CentOS/6.10
basesystem.noarch                10.0-4.el6                        @CentOS/6.10
bash.x86_64                      4.1.2-48.el6                      @CentOS/6.10
bind-libs.x86_64                 32:9.8.2-0.68.rc1.el6_10.1        @Updates/6.10
bind-utils.x86_64                32:9.8.2-0.68.rc1.el6_10.1        @Updates/6.10
binutils.x86_64                  2.20.51.0.2-5.48.el6              @CentOS/6.10
bzip2.x86_64                     1.0.5-7.el6_0                     @CentOS/6.10
bzip2-libs.x86_64                1.0.5-7.el6_0                     @CentOS/6.10
ca-certificates.noarch           2018.2.22-65.1.el6                @CentOS/6.10
centos-release.x86_64            6-10.el6.centos.12.3              @CentOS/6.10
checkpolicy.x86_64               2.0.22-1.el6                      @CentOS/6.10
chkconfig.x86_64                 1.3.49.5-1.el6                    @CentOS/6.10
coreutils.x86_64                 8.4-47.el6                        @CentOS/6.10
coreutils-libs.x86_64            8.4-47.el6                        @CentOS/6.10
cpio.x86_64                      2.10-13.el6                       @CentOS/6.10
cracklib.x86_64                  2.8.16-4.el6                      @CentOS/6.10
cracklib-dicts.x86_64            2.8.16-4.el6                      @CentOS/6.10
curl.x86_64                      7.19.7-53.el6_9                   @CentOS/6.10
cyrus-sasl-lib.x86_64            2.1.23-15.el6_6.2                 @CentOS/6.10
dash.x86_64                      0.5.5.1-4.el6                     @CentOS/6.10
db4.x86_64                       4.7.25-22.el6                     @CentOS/6.10
db4-utils.x86_64                 4.7.25-22.el6                     @CentOS/6.10
dbus-libs.x86_64                 1:1.2.24-9.el6                    @CentOS/6.10
diffutils.x86_64                 2.8.1-28.el6                      @CentOS/6.10
elfutils-libelf.x86_64           0.164-2.el6                       @CentOS/6.10
epel-release.noarch              6-8                               @extras
ethtool.x86_64                   2:3.5-6.el6                       @CentOS/6.10
expat.x86_64                     2.0.1-13.el6_8                    @CentOS/6.10
file.x86_64                      5.04-30.el6                       @CentOS/6.10
file-libs.x86_64                 5.04-30.el6                       @CentOS/6.10
filesystem.x86_64                2.4.30-3.el6                      @CentOS/6.10
findutils.x86_64                 1:4.4.2-9.el6                     @CentOS/6.10
gamin.x86_64                     0.1.10-9.el6                      @CentOS/6.10
gawk.x86_64                      3.1.7-10.el6_7.3                  @CentOS/6.10
gdbm.x86_64                      1.8.0-39.el6                      @CentOS/6.10
glib2.x86_64                     2.28.8-10.el6                     @CentOS/6.10
glibc.x86_64                     2.12-1.212.el6                    @CentOS/6.10
glibc-common.x86_64              2.12-1.212.el6                    @CentOS/6.10
gmp.x86_64                       4.3.1-13.el6                      @CentOS/6.10
gnupg2.x86_64                    2.0.14-9.el6_10                   @Updates/6.10
gpgme.x86_64                     1.1.8-3.el6                       @CentOS/6.10
grep.x86_64                      2.20-6.el6                        @CentOS/6.10
groff.x86_64                     1.18.1.4-21.el6                   @CentOS/6.10
gzip.x86_64                      1.3.12-24.el6                     @CentOS/6.10
info.x86_64                      4.13a-8.el6                       @CentOS/6.10
keyutils-libs.x86_64             1.4-5.el6                         @CentOS/6.10
krb5-libs.x86_64                 1.10.3-65.el6                     @CentOS/6.10
less.x86_64                      436-13.el6                        @CentOS/6.10
libacl.x86_64                    2.2.49-7.el6_9.1                  @CentOS/6.10
libattr.x86_64                   2.4.44-7.el6                      @CentOS/6.10
libblkid.x86_64                  2.17.2-12.28.el6_9.2              @CentOS/6.10
libcap.x86_64                    2.16-5.5.el6                      @CentOS/6.10
libcom_err.x86_64                1.41.12-24.el6                    @CentOS/6.10
libcurl.x86_64                   7.19.7-53.el6_9                   @CentOS/6.10
libffi.x86_64                    3.0.5-3.2.el6                     @CentOS/6.10
libgcc.x86_64                    4.4.7-23.el6                      @CentOS/6.10
libgcrypt.x86_64                 1.4.5-12.el6_8                    @CentOS/6.10
libgpg-error.x86_64              1.7-4.el6                         @CentOS/6.10
libidn.x86_64                    1.18-2.el6                        @CentOS/6.10
libnih.x86_64                    1.0.1-8.el6                       @CentOS/6.10
libselinux.x86_64                2.0.94-7.el6                      @CentOS/6.10
libselinux-utils.x86_64          2.0.94-7.el6                      @CentOS/6.10
libsemanage.x86_64               2.0.43-5.1.el6                    @CentOS/6.10
libsepol.x86_64                  2.0.41-4.el6                      @CentOS/6.10
libssh2.x86_64                   1.4.2-2.el6_7.1                   @CentOS/6.10
libstdc++.x86_64                 4.4.7-23.el6                      @CentOS/6.10
libtasn1.x86_64                  2.3-6.el6_5                       @CentOS/6.10
libusb.x86_64                    0.1.12-23.el6                     @CentOS/6.10
libuser.x86_64                   0.56.13-8.el6_7                   @CentOS/6.10
libutempter.x86_64               1.1.5-4.1.el6                     @CentOS/6.10
libuuid.x86_64                   2.17.2-12.28.el6_9.2              @CentOS/6.10
libxml2.x86_64                   2.7.6-21.el6_8.1                  @CentOS/6.10
lua.x86_64                       5.1.4-4.1.el6                     @CentOS/6.10
make.x86_64                      1:3.81-23.el6                     @CentOS/6.10
mingetty.x86_64                  1.08-5.el6                        @CentOS/6.10
module-init-tools.x86_64         3.9-26.el6                        @CentOS/6.10
ncurses.x86_64                   5.7-4.20090207.el6                @CentOS/6.10
ncurses-base.x86_64              5.7-4.20090207.el6                @CentOS/6.10
ncurses-libs.x86_64              5.7-4.20090207.el6                @CentOS/6.10
net-tools.x86_64                 1.60-114.el6                      @CentOS/6.10
nspr.x86_64                      4.19.0-1.el6                      @CentOS/6.10
nss.x86_64                       3.36.0-8.el6                      @CentOS/6.10
nss-softokn.x86_64               3.14.3-23.3.el6_8                 @CentOS/6.10
nss-softokn-freebl.x86_64        3.14.3-23.3.el6_8                 @CentOS/6.10
nss-sysinit.x86_64               3.36.0-8.el6                      @CentOS/6.10
nss-tools.x86_64                 3.36.0-8.el6                      @CentOS/6.10
nss-util.x86_64                  3.36.0-1.el6                      @CentOS/6.10
openldap.x86_64                  2.4.40-16.el6                     @CentOS/6.10
openssl.x86_64                   1.0.1e-57.el6                     @CentOS/6.10
p11-kit.x86_64                   0.18.5-2.el6_5.2                  @CentOS/6.10
p11-kit-trust.x86_64             0.18.5-2.el6_5.2                  @CentOS/6.10
pam.x86_64                       1.1.1-24.el6                      @CentOS/6.10
passwd.x86_64                    0.77-7.el6                        @CentOS/6.10
pcre.x86_64                      7.8-7.el6                         @CentOS/6.10
pinentry.x86_64                  0.7.6-8.el6                       @CentOS/6.10
pkgconfig.x86_64                 1:0.23-9.1.el6                    @CentOS/6.10
plymouth-core-libs.x86_64        0.8.3-29.el6.centos               @CentOS/6.10
plymouth-scripts.x86_64          0.8.3-29.el6.centos               @CentOS/6.10
popt.x86_64                      1.13-7.el6                        @CentOS/6.10
procps.x86_64                    3.2.8-45.el6_9.3                  @Updates/6.10
psmisc.x86_64                    22.6-24.el6                       @CentOS/6.10
pth.x86_64                       2.0.7-9.3.el6                     @CentOS/6.10
pygpgme.x86_64                   0.1-18.20090824bzr68.el6          @CentOS/6.10
python.x86_64                    2.6.6-66.el6_8                    @CentOS/6.10
python-iniparse.noarch           0.3.1-2.1.el6                     @CentOS/6.10
python-libs.x86_64               2.6.6-66.el6_8                    @CentOS/6.10
python-pycurl.x86_64             7.19.0-9.el6                      @CentOS/6.10
python-urlgrabber.noarch         3.9.1-11.el6                      @CentOS/6.10
readline.x86_64                  6.0-4.el6                         @CentOS/6.10
rootfiles.noarch                 8.1-6.1.el6                       @CentOS/6.10
rpm.x86_64                       4.8.0-59.el6                      @CentOS/6.10
rpm-libs.x86_64                  4.8.0-59.el6                      @CentOS/6.10
rpm-python.x86_64                4.8.0-59.el6                      @CentOS/6.10
sed.x86_64                       4.2.1-10.el6                      @CentOS/6.10
setup.noarch                     2.8.14-23.el6                     @CentOS/6.10
shadow-utils.x86_64              2:4.1.5.1-5.el6                   @CentOS/6.10
shared-mime-info.x86_64          0.70-6.el6                        @CentOS/6.10
sqlite.x86_64                    3.6.20-1.el6_7.2                  @CentOS/6.10
tar.x86_64                       2:1.23-15.el6_8                   @CentOS/6.10
tzdata.noarch                    2018e-3.el6                       @CentOS/6.10
ustr.x86_64                      1.0.4-9.1.el6                     @CentOS/6.10
vim-minimal.x86_64               2:7.4.629-5.el6_8.1               @CentOS/6.10
which.x86_64                     2.19-6.el6                        @CentOS/6.10
xz-libs.x86_64                   4.999.9-0.5.beta.20091007git.el6  @CentOS/6.10
yum.noarch                       3.2.29-81.el6.centos              @CentOS/6.10
yum-metadata-parser.x86_64       1.1.2-16.el6                      @CentOS/6.10
yum-plugin-fastestmirror.noarch  1.1.30-42.el6_10                  @Updates/6.10
yum-plugin-ovl.noarch            1.1.30-42.el6_10                  @Updates/6.10
zlib.x86_64                      1.2.3-29.el6                      @CentOS/6.10`

var dnfListArray = strings.Split(dnfList, "\n")
var yumListArray = strings.Split(yumList, "\n")

func TestParseYumList(t *testing.T) {
	var list []string
	list = append(list, "bzip2-libs.x86_64 1.0.6-13.el7")
	list = append(list, "cpio.x86_64 2.11-27.el7")
	list = append(list, "elfutils-default-yama-scope.noarch 0.172-2.el7")
	result := ParseYumList(list)

	if len(result.Projects) != 3 {
		t.Errorf("Didn't work")
	}

	if result.Projects[0].Name != "bzip2-libs" || result.Projects[0].Version != "1.0.6" {
		t.Errorf("bzip2-libs dep did not match result. Actual %s", result.Projects[0])
	}
	if result.Projects[1].Name != "cpio" || result.Projects[1].Version != "2.11" {
		t.Errorf("cpio dep did not match result. Actual %s", result.Projects[1])
	}
	if result.Projects[2].Name != "elfutils-default-yama-scope" || result.Projects[2].Version != "0.172" {
		t.Errorf("elfutils-default-yama-scope dep did not match result. Actual %s", result.Projects[2])
	}
}

func TestParseYumListFromStdIn(t *testing.T) {
	result := ParseYumListFromStdIn(yumListArray)

	if len(result.Projects) != 130 {
		t.Errorf("Didn't work")
	}

	if result.Projects[0].Name != "MAKEDEV" || result.Projects[0].Version != "3.24" {
		t.Errorf("ncurses dep did not match result. Actual %s", result.Projects[0])
	}
	if result.Projects[1].Name != "audit-libs" || result.Projects[1].Version != "2.4.5" {
		t.Errorf("coreutils dep did not match result. Actual %s", result.Projects[1])
	}
	if result.Projects[2].Name != "basesystem" || result.Projects[2].Version != "10.0" {
		t.Errorf("expat dep did not match result. Actual %s", result.Projects[2])
	}
}

func TestParseDnfListFromStdIn(t *testing.T) {
	result := ParseYumListFromStdIn(dnfListArray)

	if len(result.Projects) != 173 {
		t.Errorf("Didn't work...# projects did not match : %v", len(result.Projects))
	}

	if result.Projects[0].Name != "acl" || result.Projects[0].Version != "2.2.53" {
		t.Errorf("acl dep did not match result. Actual %s", result.Projects[0])
	}
	if result.Projects[1].Name != "audit-libs" || result.Projects[1].Version != "3.0" {
		t.Errorf("audit-libs dep did not match result. Actual %s", result.Projects[1])
	}
	if result.Projects[2].Name != "basesystem" || result.Projects[2].Version != "5" {
		t.Errorf("basesystem dep did not match result. Actual %s", result.Projects[2])
	}
}
