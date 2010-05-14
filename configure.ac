AC_INIT([pms], [0.42], [kimtjen@gmail.com])
AC_CONFIG_SRCDIR([src/pms.cpp])
MP_WITH_CURSES
HAVE_CURS=$mp_cv_curses
if test x$HAVE_CURS != xyes; then
	AC_MSG_ERROR([PMS requires the curses or ncurses library. Please install before compiling PMS.])
fi
LIBS="$LIBS $CURSES_LIB"
AC_ARG_ENABLE([regex],
       [AS_HELP_STRING([--enable-regex], [Support pattern matching through the boost regex library (default: disable)])],
       [case "${enableval}" in
               yes) AC_CHECK_LIB([boost_regex], [main]) ;;
       esac])

AC_PREFIX_PROGRAM(man)
AM_INIT_AUTOMAKE([-Wall -Werror foreign])
AC_PROG_CXX
AC_CONFIG_HEADERS([config.h])
AC_CONFIG_FILES([Makefile po/Makefile.in])

PKG_CHECK_MODULES([GLIB], [glib-2.0 >= 2.0],,
	[AC_MSG_ERROR([glib 2.0 is required])])
AM_GLIB_GNU_GETTEXT
GETTEXT_PACKAGE="pms"
AC_SUBST(GETTEXT_PACKAGE)
AC_DEFINE_UNQUOTED([GETTEXT_PACKAGE], 
                   ["${GETTEXT_PACKAGE}"], 
                   [gettext domain])
USE_NLS=$gt_cv_have_gettext
if test x$USE_NLS != xyes; then
	AC_MSG_ERROR([PMS requires gettext. Please install gettext before compiling PMS.])
fi
AC_SUBST(GLIB_CFLAGS)
AC_SUBST(GLIB_LIBS)

AC_OUTPUT
