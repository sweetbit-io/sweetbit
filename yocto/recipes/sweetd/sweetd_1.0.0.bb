SUMMARY = "sweetd"
DESCRIPTION = "Sweet daemon for pairing and control \
of the Bitcoin-enabled candy dispenser."
HOMEPAGE = "https://github.com/sweetbit-io/sweetbit"
SECTION = "net"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://git/LICENSE;md5=7087f57a125c674f2eeafee675b016a1"

SRC_URI = "\
  file://../../../../api;subdir=git;unpack=0 \
  file://../../../../app;subdir=git;unpack=0 \
  file://../../../../dispenser;subdir=git;unpack=0 \
  file://../../../../lightning;subdir=git;unpack=0 \
  file://../../../../machine;subdir=git;unpack=0 \
  file://../../../../network;subdir=git;unpack=0 \
  file://../../../../nodeman;subdir=git;unpack=0 \
  file://../../../../onion;subdir=git;unpack=0 \
  file://../../../../pairing;subdir=git;unpack=0 \
  file://../../../../pos;subdir=git;unpack=0 \
  file://../../../../reboot;subdir=git;unpack=0 \
  file://../../../../state;subdir=git;unpack=0 \
  file://../../../../sweetdb;subdir=git;unpack=0 \
  file://../../../../sweetlog;subdir=git;unpack=0 \
  file://../../../../sysid;subdir=git;unpack=0 \
  file://../../../../updater;subdir=git;unpack=0 \
  file://../../../../config.go;subdir=git;unpack=0 \
  file://../../../../go.mod;subdir=git;unpack=0 \
  file://../../../../go.sum;subdir=git;unpack=0 \
  file://../../../../LICENSE;subdir=git;unpack=0 \
  file://../../../../main.go;subdir=git;unpack=0 \
  file://../../../../Makefile;subdir=git;unpack=0 \
  file://sweetd.init \
  file://sweetd.default \
  file://sweet.conf \
  file://sweetd.service \
  "

PROVIDES = "sweetd"
RPROVIDES_${PN} = "sweetd"

DEPENDS += "packr2-native"
RDEPENDS_${PN} += "wpa-supplicant bluez5 pi-bluetooth lnd tor"

FILESEXTRAPATHS_prepend := "${THISDIR}/files:"

inherit update-rc.d systemd npm go

INITSCRIPT_NAME = "sweetd"
INITSCRIPT_PARAMS = "defaults 92 20"
INSANE_SKIP_${PN} = "already-stripped"

INSANE_SKIP_${PN} = "ldflags"
INSANE_SKIP_${PN}-dev = "ldflags"

S = "${WORKDIR}"

python do_unpack() {
  bb.build.exec_func('base_do_unpack', d)
}

do_configure() {
  base_do_configure
}
do_configure[deptask] = "do_populate_staging"

do_compile() {
  cd ${S}/git
  oe_runmake
}

do_install() {
  install -m 0755 -d ${D}${bindir} ${D}${docdir}/sweetd
  install -m 0755 ${S}/git/sweetd ${D}${bindir}/sweetd

  install -m 0755 -d ${D}${sysconfdir}/init.d ${D}${sysconfdir}/default ${D}${sysconfdir}/dbus-1 ${D}${sysconfdir}/dbus-1/system.d
  install -m 0755 ${WORKDIR}/sweetd.init ${D}${sysconfdir}/init.d/sweetd
  sed -i -e 's,@BINDIR@,${bindir},g' -e 's,@SYSCONFDIR@,${sysconfdir},g' ${D}${sysconfdir}/init.d/sweetd
  install -m 0755 ${WORKDIR}/sweetd.default ${D}${sysconfdir}/default/sweetd
  install -m 0755 ${WORKDIR}/sweet.conf ${D}${sysconfdir}/dbus-1/system.d/

  install -m 0755 -d ${D}${systemd_unitdir}/system/
  install -m 0644 ${WORKDIR}/sweetd.service ${D}${systemd_unitdir}/system/
  sed -i -e 's,@BINDIR@,${bindir},g' -e 's,@SYSCONFDIR@,${sysconfdir},g' ${D}${systemd_unitdir}/system/sweetd.service
}

FILES_${PN} = "\
  ${bindir}/sweetd \
  ${sysconfdir}/init.d/sweetd \
  ${sysconfdir}/default/sweetd \
  ${sysconfdir}/dbus-1/system.d/sweet.conf \
  ${systemd_unitdir}/system/sweetd.service \
  "
