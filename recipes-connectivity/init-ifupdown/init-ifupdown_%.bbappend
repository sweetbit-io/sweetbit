SRC_URI += "\
    file://interfaces \
    "

FILESEXTRAPATHS_prepend := "${THISDIR}/files:"

do_install_append() {
     install -m 0644 ${WORKDIR}/interfaces ${D}${sysconfdir}/network/interfaces
}
