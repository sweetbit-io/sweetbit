SUMMARY = "Lightning Network Daemon"
DESCRIPTION = "Lightning is a decentralized network using \
smart contract functionality in the blockchain to enable \
instant payments across a network of participants."
HOMEPAGE = "https://github.com/lightningnetwork/lnd"
SECTION = "net"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://src/${GO_IMPORT}/LICENSE;md5=93518a723211417febbbbd1c5230f83f"

GO_IMPORT = "github.com/lightningnetwork/lnd"
SRC_URI = " \
  git://${GO_IMPORT} \
  "

FILESEXTRAPATHS_prepend := "${THISDIR}/files:"
SRC_URI += " file://0001-Enable-neutrino-on-mainnet.patch"

PROVIDES = "lnd"
RPROVIDES_${PN} = "lnd"

SRCREV = "add905d17f7bbb11d0df2761cdf8accf2fef2b00"

inherit go

GO_INSTALL = "${GO_IMPORT}/cmd/lnd"

do_compile () {
  cd ${S}/src/${GO_IMPORT}
  oe_runmake install
}

do_install () {
	install -m 0755 -d ${D}${bindir}
  install -m 0755 ${B}/bin/${GOOS}_${GOARCH}/lnd ${D}${bindir}
	install -m 0755 ${B}/bin/${GOOS}_${GOARCH}/lncli ${D}${bindir}
}

FILES_${PN} = "\
  ${bindir}/lnd \
  ${bindir}/lncli \
  "
