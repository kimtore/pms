bin_PROGRAMS = pms
man_MANS = pms.1
EXTRA_pms_SOURCES = src/action.h src/color.h src/command.h src/config.h src/conn.h src/display.h src/input.h src/list.h src/song.h src/types.h src/libmpdclient.h src/i18n.h src/topbar.h src/pms.h src/settings.h src/field.h src/mycurses.h src/mediator.h src/message.h src/filter.h
EXTRA_DIST = config.rpath m4/ChangeLog  pms.1 po
pms_SOURCES = src/action.cpp src/color.cpp src/command.cpp src/config.cpp src/conn.cpp src/display.cpp src/input.cpp src/list.cpp src/pms.cpp src/song.cpp src/libmpdclient.c src/field.cpp src/settings.cpp src/mediator.cpp src/message.cpp
CXXFLAGS += \
	@GLIB_CFLAGS@ \
	-DLOCALE_DIR=\""$(datadir)/locale"\"

pms_LDADD = \
	$(GLIB_LIBS)

ACLOCAL_AMFLAGS = -I m4
