FILESEXTRAPATHS_prepend := "${THISDIR}/files:"

# Add Red Hat NTP servers to synchronize the UDOO local date
SRC_URI += "file://ntp.conf"

do_install_append() {
  install -m 644 ${WORKDIR}/ntp.conf ${D}${sysconfdir}
}
