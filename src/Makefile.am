bin_PROGRAMS = pms
EXTRA_pms_SOURCES = *.h
pms_SOURCES = action.cpp color.cpp command.cpp config.cpp conn.cpp display.cpp input.cpp list.cpp listitem.cpp songlist.cpp pms.cpp song.cpp field.cpp options.cpp message.cpp playlist.cpp set_parameters.cpp

AM_CXXFLAGS = \
	@glib_CFLAGS@ \
	@libmpdclient_CFLAGS@ \
	-DLOCALE_DIR=\""$(datadir)/locale"\"

pms_LDADD = @CURSES_LIB@ $(glib_LIBS) $(libmpdclient_LIBS) -lpthread

ACLOCAL_AMFLAGS = -I m4
