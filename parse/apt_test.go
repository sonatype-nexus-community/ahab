//
// Copyright (c) 2019-present Sonatype, Inc.
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
//

package parse_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/sonatype-nexus-community/ahab/parse"
)

func TestParseAptList(t *testing.T) {
	var list []string
	list = append(list, "libedit2 3.1-20170329-1")
	list = append(list, "libmount1 2.31.1-0.4ubuntu3.3")
	list = append(list, "zlib1g 1:1.2.11.dfsg-0ubuntu2")
	result := ParseAptList(list)

	if len(result.Projects) != 3 {
		t.Errorf("Didn't work")
	}
}

// generate CLI package list via:
// # dpkg-query --show --showformat='${Package} ${Version}\n'
var dpkgList = `adduser 3.116ubuntu1
apt 1.6.12
base-files 10.1ubuntu2.6
base-passwd 3.5.44
bash 4.4.18-2ubuntu1.2
binutils 2.30-21ubuntu1~18.04.2
binutils-common 2.30-21ubuntu1~18.04.2
binutils-x86-64-linux-gnu 2.30-21ubuntu1~18.04.2
bsdutils 1:2.31.1-0.4ubuntu3.4
build-essential 12.4ubuntu1
bzip2 1.0.6-8.1ubuntu0.2
ca-certificates 20180409
coreutils 8.28-1ubuntu1
cpp 4:7.4.0-1ubuntu2.3
cpp-7 7.4.0-1ubuntu1~18.04.1
cron 3.0pl1-128.1ubuntu1
dash 0.5.8-2.10
dbus 1.12.2-1ubuntu1.1
debconf 1.5.66ubuntu1
debianutils 4.8.4
diffutils 1:3.6-1
dirmngr 2.2.4-1ubuntu1.2
distro-info-data 0.37ubuntu0.5
dpkg 1.19.0.5ubuntu2.2
dpkg-dev 1.19.0.5ubuntu2.2
e2fsprogs 1.44.1-1ubuntu1.1
fakeroot 1.22-2ubuntu1
fdisk 2.31.1-0.4ubuntu3.4
file 1:5.32-2ubuntu0.2
findutils 4.6.0+git+20170828-2
g++ 4:7.4.0-1ubuntu2.3
g++-7 7.4.0-1ubuntu1~18.04.1
gcc 4:7.4.0-1ubuntu2.3
gcc-7 7.4.0-1ubuntu1~18.04.1
gcc-7-base 7.4.0-1ubuntu1~18.04.1
gcc-8-base 8.3.0-6ubuntu1~18.04.1
gir1.2-glib-2.0 1.56.1-1
gnupg 2.2.4-1ubuntu1.2
gnupg-l10n 2.2.4-1ubuntu1.2
gnupg-utils 2.2.4-1ubuntu1.2
golang-1.13-go 1.13-1longsleep2+xenial
golang-1.13-src 1.13-1longsleep2+xenial
golang-go 2:1.13~1longsleep1+xenial
golang-src 2:1.13~1longsleep1+xenial
gpg 2.2.4-1ubuntu1.2
gpg-agent 2.2.4-1ubuntu1.2
gpg-wks-client 2.2.4-1ubuntu1.2
gpg-wks-server 2.2.4-1ubuntu1.2
gpgconf 2.2.4-1ubuntu1.2
gpgsm 2.2.4-1ubuntu1.2
gpgv 2.2.4-1ubuntu1.2
grep 3.1-2
gzip 1.6-5ubuntu1
hostname 3.20
init-system-helpers 1.51
iso-codes 3.79-1
libacl1 2.2.52-3build1
libalgorithm-diff-perl 1.19.03-1
libalgorithm-diff-xs-perl 0.04-5
libalgorithm-merge-perl 0.08-3
libapparmor1 2.12-4ubuntu5.1
libapt-inst2.0 1.6.12
libapt-pkg5.0 1.6.12
libasan4 7.4.0-1ubuntu1~18.04.1
libasn1-8-heimdal 7.5.0+dfsg-1
libassuan0 2.5.1-2
libatomic1 8.3.0-6ubuntu1~18.04.1
libattr1 1:2.4.47-2build1
libaudit-common 1:2.8.2-1ubuntu1
libaudit1 1:2.8.2-1ubuntu1
libbinutils 2.30-21ubuntu1~18.04.2
libblkid1 2.31.1-0.4ubuntu3.4
libbz2-1.0 1.0.6-8.1ubuntu0.2
libc-bin 2.27-3ubuntu1
libc-dev-bin 2.27-3ubuntu1
libc6 2.27-3ubuntu1
libc6-dev 2.27-3ubuntu1
libcap-ng0 0.7.7-3.1
libcc1-0 8.3.0-6ubuntu1~18.04.1
libcilkrts5 7.4.0-1ubuntu1~18.04.1
libcom-err2 1.44.1-1ubuntu1.1
libdb5.3 5.3.28-13.1ubuntu1.1
libdbus-1-3 1.12.2-1ubuntu1.1
libdebconfclient0 0.213ubuntu1
libdpkg-perl 1.19.0.5ubuntu2.2
libexpat1 2.2.5-3ubuntu0.2
libext2fs2 1.44.1-1ubuntu1.1
libfakeroot 1.22-2ubuntu1
libfdisk1 2.31.1-0.4ubuntu3.4
libffi6 3.2.1-8
libfile-fcntllock-perl 0.22-3build2
libgcc-7-dev 7.4.0-1ubuntu1~18.04.1
libgcc1 1:8.3.0-6ubuntu1~18.04.1
libgcrypt20 1.8.1-4ubuntu1.1
libgdbm-compat4 1.14.1-6
libgdbm5 1.14.1-6
libgirepository-1.0-1 1.56.1-1
libglib2.0-0 2.56.4-0ubuntu0.18.04.4
libglib2.0-data 2.56.4-0ubuntu0.18.04.4
libgmp10 2:6.1.2+dfsg-2
libgnutls30 3.5.18-1ubuntu1.1
libgomp1 8.3.0-6ubuntu1~18.04.1
libgpg-error0 1.27-6
libgpm2 1.20.7-5
libgssapi3-heimdal 7.5.0+dfsg-1
libhcrypto4-heimdal 7.5.0+dfsg-1
libheimbase1-heimdal 7.5.0+dfsg-1
libheimntlm0-heimdal 7.5.0+dfsg-1
libhogweed4 3.4-1
libhx509-5-heimdal 7.5.0+dfsg-1
libicu60 60.2-3ubuntu3
libidn2-0 2.0.4-1.1build2
libisl19 0.19-1
libitm1 8.3.0-6ubuntu1~18.04.1
libkrb5-26-heimdal 7.5.0+dfsg-1
libksba8 1.3.5-2
libldap-2.4-2 2.4.45+dfsg-1ubuntu1.4
libldap-common 2.4.45+dfsg-1ubuntu1.4
liblocale-gettext-perl 1.07-3build2
liblsan0 8.3.0-6ubuntu1~18.04.1
liblz4-1 0.0~r131-2ubuntu3
liblzma5 5.2.2-1.3
libmagic-mgc 1:5.32-2ubuntu0.2
libmagic1 1:5.32-2ubuntu0.2
libmount1 2.31.1-0.4ubuntu3.4
libmpc3 1.1.0-1
libmpdec2 2.4.2-1ubuntu1
libmpfr6 4.0.1-1
libmpx2 8.3.0-6ubuntu1~18.04.1
libncurses5 6.1-1ubuntu1.18.04
libncursesw5 6.1-1ubuntu1.18.04
libnettle6 3.4-1
libnpth0 1.5-3
libp11-kit0 0.23.9-2
libpam-modules 1.1.8-3.6ubuntu2.18.04.1
libpam-modules-bin 1.1.8-3.6ubuntu2.18.04.1
libpam-runtime 1.1.8-3.6ubuntu2.18.04.1
libpam0g 1.1.8-3.6ubuntu2.18.04.1
libpcre3 2:8.39-9
libperl5.26 5.26.1-6ubuntu0.3
libprocps6 2:3.3.12-3ubuntu1.2
libpython3-stdlib 3.6.7-1~18.04
libpython3.6 3.6.8-1~18.04.2
libpython3.6-minimal 3.6.8-1~18.04.2
libpython3.6-stdlib 3.6.8-1~18.04.2
libquadmath0 8.3.0-6ubuntu1~18.04.1
libreadline7 7.0-3
libroken18-heimdal 7.5.0+dfsg-1
libsasl2-2 2.1.27~101-g0780600+dfsg-3ubuntu2
libsasl2-modules 2.1.27~101-g0780600+dfsg-3ubuntu2
libsasl2-modules-db 2.1.27~101-g0780600+dfsg-3ubuntu2
libseccomp2 2.4.1-0ubuntu0.18.04.2
libselinux1 2.7-2build2
libsemanage-common 2.7-2build2
libsemanage1 2.7-2build2
libsepol1 2.7-1
libsmartcols1 2.31.1-0.4ubuntu3.4
libsqlite3-0 3.22.0-1ubuntu0.1
libss2 1.44.1-1ubuntu1.1
libssl1.1 1.1.1-1ubuntu2.1~18.04.4
libstdc++-7-dev 7.4.0-1ubuntu1~18.04.1
libstdc++6 8.3.0-6ubuntu1~18.04.1
libsystemd0 237-3ubuntu10.29
libtasn1-6 4.13-2
libtinfo5 6.1-1ubuntu1.18.04
libtsan0 8.3.0-6ubuntu1~18.04.1
libubsan0 7.4.0-1ubuntu1~18.04.1
libudev1 237-3ubuntu10.29
libunistring2 0.9.9-0ubuntu2
libuuid1 2.31.1-0.4ubuntu3.4
libwind0-heimdal 7.5.0+dfsg-1
libxml2 2.9.4+dfsg1-6.1ubuntu1.2
libzstd1 1.3.3+dfsg-2ubuntu1.1
linux-libc-dev 4.15.0-64.73
login 1:4.5-1ubuntu2
lsb-base 9.20170808ubuntu1
lsb-release 9.20170808ubuntu1
make 4.1-9.1ubuntu1
manpages 4.15-1
manpages-dev 4.15-1
mawk 1.3.3-17ubuntu3
mime-support 3.60ubuntu1
mount 2.31.1-0.4ubuntu3.4
ncurses-base 6.1-1ubuntu1.18.04
ncurses-bin 6.1-1ubuntu1.18.04
netbase 5.4
openssl 1.1.1-1ubuntu2.1~18.04.4
passwd 1:4.5-1ubuntu2
patch 2.7.6-2ubuntu1.1
perl 5.26.1-6ubuntu0.3
perl-base 5.26.1-6ubuntu0.3
perl-modules-5.26 5.26.1-6ubuntu0.3
pinentry-curses 1.1.0-1
pkg-config 0.29.1-0ubuntu2
powermgmt-base 1.33
procps 2:3.3.12-3ubuntu1.2
python-apt-common 1.6.4
python3 3.6.7-1~18.04
python3-apt 1.6.4
python3-dbus 1.2.6-1
python3-gi 3.26.1-2ubuntu1
python3-minimal 3.6.7-1~18.04
python3-software-properties 0.96.24.32.11
python3.6 3.6.8-1~18.04.2
python3.6-minimal 3.6.8-1~18.04.2
readline-common 7.0-3
sed 4.4-2
sensible-utils 0.0.12
shared-mime-info 1.9-2
software-properties-common 0.96.24.32.11
sysvinit-utils 2.88dsf-59.10ubuntu1
tar 1.29b-2ubuntu0.1
ubuntu-keyring 2018.09.18.1~18.04.0
ucf 3.0038
unattended-upgrades 1.1ubuntu1.18.04.11
util-linux 2.31.1-0.4ubuntu3.4
vim 2:8.0.1453-1ubuntu1.1
vim-common 2:8.0.1453-1ubuntu1.1
vim-runtime 2:8.0.1453-1ubuntu1.1
xdg-user-dirs 0.17-1ubuntu1
xxd 2:8.0.1453-1ubuntu1.1
xz-utils 5.2.2-1.3
zlib1g 1:1.2.11.dfsg-0ubuntu2

`

var dpkgListArray = strings.Split(dpkgList, "\n")

func TestParseDpkgList(t *testing.T) {
	result := ParseDpkgList(dpkgListArray)

	// adduser 3.116ubuntu1
	assert.Equal(t, Projects{"adduser", "3.116"}, result.Projects[0])

	// apt 1.6.12
	assert.Equal(t, Projects{"apt", "1.6.12"}, result.Projects[1])

	// ca-certificates 20180409
	assert.Equal(t, Projects{"ca-certificates", "20180409"}, result.Projects[11])

	// @todo Is the resulting version correct for this case?
	// diffutils 1:3.6-1
	assert.Equal(t, Projects{"diffutils", "3.6"}, result.Projects[20])

	// libsystemd0 237-3ubuntu10.29
	assert.Equal(t, Projects{"libsystemd0", "237-3"}, result.Projects[162])

	// libudev1 237-3ubuntu10.29
	assert.Equal(t, Projects{"libudev1", "237-3"}, result.Projects[167])

	// tar 1.29b-2ubuntu0.1
	assert.Equal(t, Projects{"tar", "1.29"}, result.Projects[211])

	// vim 2:8.0.1453-1ubuntu1.1
	assert.Equal(t, Projects{"vim", "8.0.1453"}, result.Projects[216])

	// @todo Is the resulting version correct for this case?
	// xz-utils 5.2.2-1.3
	assert.Equal(t, Projects{"xz-utils", "5.2.2"}, result.Projects[221])

	// zlib1g 1:1.2.11.dfsg-0ubuntu2
	assert.Equal(t, Projects{"zlib1g", "1.2.11"}, result.Projects[222])
}
