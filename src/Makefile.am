bin_PROGRAMS = pms
EXTRA_pms_SOURCES = *.h
pms_SOURCES = \
    action.cpp \
    color.cpp \
    command.cpp \
    config.cpp \
    conn.cpp \
    display.cpp \
    error.cpp \
    field.cpp \
    input.cpp \
    library.cpp \
    list.cpp \
    listitem.cpp \
    message.cpp \
    options.cpp \
    playlist.cpp \
    pms.cpp \
    queue.cpp \
	search.cpp \
    set_parameters.cpp \
    song.cpp \
    songlist.cpp

AM_CXXFLAGS = \
	@glib_CFLAGS@ \
	@libmpdclient_CFLAGS@ \
	-DLOCALE_DIR=\""$(datadir)/locale"\"

pms_LDADD = @CURSES_LIB@ $(glib_LIBS) $(libmpdclient_LIBS) -lpthread

ACLOCAL_AMFLAGS = -I m4
