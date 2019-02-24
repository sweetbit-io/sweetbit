DESCRIPTION = "Sweet daemon for pairing and control of the Bitcoin-enabled candy dispenser"
HOMEPAGE = "https://the.lightning.land"

VERSION = "0.4.7"
ARCHIVE = ""
SHA256SUM = "c8f7c4e1c25fb2b04fdf6b756526c66c39d52b55a547869c4955ed96ac0bb5ed"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${WORKDIR}/sweetd_${VERSION}_linux_armv6/LICENSE;md5=7087f57a125c674f2eeafee675b016a1"

PR = "r0"
PROVIDES = "sweetd"
RPROVIDES_${PN} = "sweetd"

RDEPENDS_${PN} = " wpa-supplicant iw hostapd dnsmasq"

SRC_URI = "\
    https://github.com/the-lightning-land/sweetd/releases/download/v${VERSION}/sweetd_${VERSION}_linux_armv6.tar.gz;sha256sum=${SHA256SUM} \
    file://init \
    file://default \
    file://sweetd.service \
    "

inherit update-rc.d systemd

INITSCRIPT_NAME = "sweetd"
INITSCRIPT_PARAMS = "defaults 92 20"
INSANE_SKIP_${PN} = "already-stripped"

do_configure() {
}

do_install() {
    install -m 0755 -d ${D}${bindir} ${D}${docdir}/sweetd
    install -m 0755 ${WORKDIR}/sweetd_${VERSION}_linux_armv6/sweetd ${D}${bindir}/sweetd
    install -m 0644 ${WORKDIR}/sweetd_${VERSION}_linux_armv6/README.md ${D}${docdir}/sweetd/README

    install -m 0755 -d ${D}${sysconfdir}/init.d ${D}${sysconfdir}/default
    install -m 0755 ${WORKDIR}/init ${D}${sysconfdir}/init.d/sweetd
    install -m 0755 ${WORKDIR}/default ${D}${sysconfdir}/default/sweetd
    sed -i -e 's,@BINDIR@,${bindir},g' -e 's,@SYSCONFDIR@,${sysconfdir},g' ${D}${sysconfdir}/init.d/sweetd

    install -m 0755 -d ${D}${systemd_unitdir}/system/
    install -m 0644 ${WORKDIR}/sweetd.service ${D}${systemd_unitdir}/system/
    sed -i -e 's,@BINDIR@,${bindir},g' -e 's,@SYSCONFDIR@,${sysconfdir},g' ${D}${systemd_unitdir}/system/sweetd.service
}

FILES_${PN} = "\
    ${bindir}/sweetd \
    ${sysconfdir}/init.d/sweetd \
    ${sysconfdir}/default/sweetd \
    ${systemd_unitdir}/system/sweetd.service \
    "
