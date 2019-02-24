DESCRIPTION = "Sweet daemon for pairing and control of the Bitcoin-enabled candy dispenser"
HOMEPAGE = "https://the.lightning.land"

LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${WORKDIR}/sweetd_0.4.6_linux_armv6/LICENSE;md5=7087f57a125c674f2eeafee675b016a1"

PR = "r0"
PROVIDES = "sweetd"
RPROVIDES_${PN} = "sweetd"

RDEPENDS_${PN} = " wpa-supplicant iw hostapd dnsmasq"

SRC_URI = "https://github.com/the-lightning-land/sweetd/releases/download/v0.4.6/sweetd_0.4.6_linux_armv6.tar.gz;sha256sum=69a748b7ed599075d02661117dc76e60633df42595544d18899d1bc86c21f6b6 \
    file://init;name=init \
    file://default;name=default"

inherit update-rc.d systemd

INITSCRIPT_NAME = "sweetd"
INITSCRIPT_PARAMS="defaults 40"

INSANE_SKIP_${PN} = "already-stripped"

do_configure() {
}

do_install() {
    install -m 0755 -d ${D}${bindir} ${D}${docdir}/sweetd
    install -m 0755 ${WORKDIR}/sweetd_0.4.6_linux_armv6/sweetd ${D}${bindir}/sweetd
    install -m 0644 ${WORKDIR}/sweetd_0.4.6_linux_armv6/README.md ${D}${docdir}/sweetd/README

    install -m 0755 -d ${D}${sysconfdir}/init.d ${D}${sysconfdir}/default
    install -m 0755 ${WORKDIR}/init ${D}${sysconfdir}/init.d/sweetd
    install -m 0755 ${WORKDIR}/default ${D}${sysconfdir}/default/sweetd
}

FILES_${PN} = "${bindir}/sweetd ${sysconfdir}/init.d/sweetd ${sysconfdir}/default/sweetd"
