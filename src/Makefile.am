bin_PROGRAMS = pms
EXTRA_pms_SOURCES = action.h color.h command.h config.h conn.h display.h input.h list.h song.h types.h i18n.h topbar.h pms.h settings.h field.h mycurses.h mediator.h message.h filter.h playlist.h
pms_SOURCES = action.cpp color.cpp command.cpp config.cpp conn.cpp display.cpp input.cpp list.cpp pms.cpp song.cpp field.cpp settings.cpp mediator.cpp message.cpp playlist.cpp

AM_CXXFLAGS = \
	@glib_CFLAGS@ \
	@libzmq_CFLAGS@ \
	@libmpdclient_CFLAGS@ \
	-DLOCALE_DIR=\""$(datadir)/locale"\"

pms_LDADD = @CURSES_LIB@ $(glib_LIBS) $(libzmq_LIBS) $(libmpdclient_LIBS) -lpthread

ACLOCAL_AMFLAGS = -I m4
