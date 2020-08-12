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


// generate CLI package list via:
// # yum list installed
var dpkgList = `Loaded plugins: fastestmirror, ovl
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

var yumListArray = strings.Split(dpkgList, "\n")

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
