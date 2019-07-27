SUMMARY = "The Tor anonymity network"
DESCRIPTION = "Tor is free software and an open network \
that helps you defend against traffic analysis, \
a form of network surveillance that threatens personal freedom and privacy, \
confidential business activities and relationships, and state security."
HOMEPAGE = "https://www.torproject.org"
SECTION = "net"
LICENSE = "BSD-3-Clause"
LIC_FILES_CHKSUM = "file://LICENSE;md5=5f5073beceebf4a374e5fae87ff912b2"

SRC_URI = "https://dist.torproject.org/tor-${PV}.tar.gz"
SRC_URI[md5sum] = "23278fc58d0014db22b428cdae3ea966"
SRC_URI[sha256sum] = "b5a2cbf0dcd3f1df2675dbd5ec10bbe6f8ae995c41b68cebe2bc95bffc90696e"

DEPENDS = "xz libevent openssl"

inherit python3native pythonnative perlnative pkgconfig autotools

EXTRA_OECONF = "--disable-tool-name-check"
