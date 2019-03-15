DESCRIPTION = "Recipe for adding persistent.txt to /data"
LICENSE = "MIT"

FILES_${PN} += "/data/persistent.txt"

do_compile() {
    echo 'This partition is persistent.' > persistent.txt
}

do_install() {
    install -d ${D}/data/
    install -m 0644 persistent.txt ${D}/data/
}