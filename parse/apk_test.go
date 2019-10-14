package parse_test

import (
	"strings"
	"testing"

	. "github.com/sonatype-nexus-community/ahab/parse"
	"github.com/sonatype-nexus-community/nancy/types"
	"github.com/stretchr/testify/assert"
)

func TestParseApkShow(t *testing.T) {
	var list []string
	list = append(list, "alpine-baselayout-3.1.2-r0 - Alpine base dir structure and init scripts")
	list = append(list, "alpine-keys-2.1-r2 - Public keys for Alpine Linux packages")
	list = append(list, "apk-tools-2.10.4-r2 - Alpine Package Keeper - package manager for alpine")
	result := ParseApkShow(list)

	if len(result.Projects) != 3 {
		t.Errorf("Didn't work")
	}
}

// generate CLI package list via:
// # apk info -vv | sort
var apkShowList = `WARNING: Ignoring APKINDEX.00740ba1.tar.gz: No such file or directory
WARNING: Ignoring APKINDEX.d8b2a6f4.tar.gz: No such file or directory
alpine-baselayout-3.1.2-r0 - Alpine base dir structure and init scripts
alpine-keys-2.1-r2 - Public keys for Alpine Linux packages
apk-tools-2.10.4-r2 - Alpine Package Keeper - package manager for alpine
busybox-1.30.1-r2 - Size optimized toolbox of many common UNIX utilities
ca-certificates-cacert-20190108-r0 - Mozilla bundled certificates
libc-utils-0.7.1-r0 - Meta package to pull in correct libc
libcrypto1.1-1.1.1c-r0 - Crypto library from openssl
libssl1.1-1.1.1c-r0 - SSL shared libraries
libtls-standalone-2.9.1-r0 - libtls extricated from libressl sources
musl-1.1.22-r3 - the musl c library (libc) implementation
musl-utils-1.1.22-r3 - the musl c library (libc) implementation
scanelf-1.2.3-r0 - Scan ELF binaries for stuff
ssl_client-1.30.1-r2 - EXternal ssl_client for busybox wget
zlib-1.2.11-r1 - A compression/decompression Library`

var apkShowArray = strings.Split(apkShowList, "\n")

func TestParseApkShowList(t *testing.T) {
	result := ParseApkShow(apkShowArray)

	// adduser 3.116ubuntu1
	assert.Equal(t, types.Projects{"alpine-baselayout", "3.1.2-r0"}, result.Projects[0])

	// apt 1.6.12
	assert.Equal(t, types.Projects{"alpine-keys", "2.1-r2"}, result.Projects[1])

	// ca-certificates 20180409
	assert.Equal(t, types.Projects{"apk-tools", "2.10.4-r2"}, result.Projects[2])
}
