SRC_URI += "\
    file://interfaces \
    "

FILESEXTRAPATHS_prepend := "${THISDIR}/files:"

# Override default so it starts second, after dbus
INITSCRIPT_PARAMS = "start 02 2 3 4 5 . stop 80 0 6 1 ."

do_install_append() {
     install -m 0644 ${WORKDIR}/interfaces ${D}${sysconfdir}/network/interfaces
}
