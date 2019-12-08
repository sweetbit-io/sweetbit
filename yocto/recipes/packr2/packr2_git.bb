SUMMARY = "packr2"
DESCRIPTION = "The simple and easy way to embed static files into Go binaries."
HOMEPAGE = "https://github.com/gobuffalo/packr"
SECTION = "net"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://git/LICENSE.txt;md5=85a1cdcf71593cf8d843001b33ab4265"

SRC_URI = "git://github.com/gobuffalo/packr"
SRCREV = "662c20c19dde9677ffbd6f107aa8d68538a3ae95"
PV = "${SRCREV}"

PROVIDES = "packr2"
RPROVIDES_${PN} = "packr2"

BBCLASSEXTEND = "native nativesdk"

inherit go

S = "${WORKDIR}"

python do_unpack() {
  bb.build.exec_func('base_do_unpack', d)
}

do_configure() {
  base_do_configure
}

do_compile() {
  cd ${S}/git
  go build ${GOBUILDFLAGS} -o packr2 github.com/gobuffalo/packr/v2/packr2
}

do_install() {
  install -m 0755 -d ${D}${bindir}
  install -m 0755 ${S}/git/packr2 ${D}${bindir}/packr2
}

FILES_${PN} = "\
  ${bindir}/packr2 \
  "
