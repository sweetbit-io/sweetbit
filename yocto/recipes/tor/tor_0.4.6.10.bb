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
SRC_URI[md5sum] = "1da676163e4c78efcc650210fa7c0530"
SRC_URI[sha256sum] = "94ccd60e04e558f33be73032bc84ea241660f92f58cfb88789bda6893739e31c"

DEPENDS = "xz libevent openssl"

inherit python3native perlnative pkgconfig autotools

EXTRA_OECONF = "--disable-tool-name-check"
