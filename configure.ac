#                                               -*- Autoconf -*-
# Process this file with autoconf to produce a configure script.

AC_PREREQ([2.69])
AC_INIT([pms], [0.99.0~git], [kimtjen@gmail.com])
AM_INIT_AUTOMAKE([-Wall -Werror foreign])
AC_CONFIG_SRCDIR([src/pms.cpp])
AC_CONFIG_HEADERS([config.h])
AC_CONFIG_MACRO_DIR([m4])
AC_LANG(C++)
# Enable debugging?
AX_CHECK_ENABLE_DEBUG

# Checks for programs.
AC_PROG_AWK
AC_PROG_CXX
AC_PROG_CC
AC_PROG_INSTALL
AC_PROG_LN_S
AC_PROG_MAKE_SET

# Checks for libraries.
AX_WITH_CURSES
PKG_CHECK_MODULES([glib], [glib-2.0 >= 2.0],, [AC_MSG_ERROR([glib 2.0 or newer is required])])
PKG_CHECK_MODULES([libmpdclient], [libmpdclient >= 2.5],, [AC_MSG_ERROR([libmpdclient 2.5 or newer is required])])

# Check for Pandoc
AC_PATH_PROG(PANDOC,pandoc,no)
if [[ "x$PANDOC" = xno ]]; then
    AC_MSG_WARN([The manual page will not be built.])
fi
AM_CONDITIONAL([PANDOC], [test "x$PANDOC" = xyes])

# Checks for header files.
AC_CHECK_HEADERS([arpa/inet.h fcntl.h libintl.h limits.h locale.h netdb.h netinet/in.h stdlib.h string.h sys/param.h sys/socket.h sys/time.h unistd.h])

# Checks for typedefs, structures, and compiler characteristics.
AC_CHECK_HEADER_STDBOOL
AC_TYPE_SIZE_T

# Checks for library functions.
AC_FUNC_MALLOC
AC_FUNC_REALLOC
AC_CHECK_FUNCS([gethostbyname memmove memset select setlocale socket strchr strdup strerror strstr strtol])

# Gettext
IT_PROG_INTLTOOL
AM_GNU_GETTEXT([external])
GETTEXT_PACKAGE="pms"
AC_SUBST([GETTEXT_PACKAGE])
AC_DEFINE_UNQUOTED([GETTEXT_PACKAGE], 
                   ["${GETTEXT_PACKAGE}"], 
                   [gettext domain])

# Other stuff
AC_ARG_ENABLE([regex],
              [AS_HELP_STRING([--enable-regex], [Support pattern matching through c++11 regex (needs gcc >= 4.9 or other compatible compiler) (default: disable)])],
       [case "${enableval}" in
            yes) AX_CXX_COMPILE_STDCXX_11([noext], [mandatory])
            AC_DEFINE([HAVE_REGEX], [1], [C++ 11 is available and regex is part of it]);;
       esac])

AC_CONFIG_FILES([Makefile
                 po/Makefile.in
                 src/Makefile
                 doc/Makefile])
AC_OUTPUT
