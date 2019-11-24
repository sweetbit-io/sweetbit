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
SRC_URI[md5sum] = "7a7b414dca81b87d3c51569fe5dce376"
SRC_URI[sha256sum] = "1e92b655a19062209c747c2f328f2b52009d8256a9514705bb8a6cfebb21b3ae"

DEPENDS = "xz libevent openssl"

inherit python3native pythonnative perlnative pkgconfig autotools

EXTRA_OECONF = "--disable-tool-name-check"
