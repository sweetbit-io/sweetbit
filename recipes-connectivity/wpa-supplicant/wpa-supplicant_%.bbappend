FILESEXTRAPATHS_prepend := "${THISDIR}/files:"

SRC_URI += " file://wpa_supplicant.conf"

FILES_${PN} += "${sysconfdir}/wpa_supplicant.conf"

do_install_append() {
  install -m 0755 ${WORKDIR}/wpa_supplicant.conf ${D}${sysconfdir}/wpa_supplicant.conf
}
