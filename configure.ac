#                                               -*- Autoconf -*-
# Process this file with autoconf to produce a configure script.

AC_PREREQ([2.68])
AC_INIT(pms, 0.50, kimtjen@gmail.com)
AM_INIT_AUTOMAKE([-Wall -Werror foreign])
AC_COPYRIGHT(Copyright (c) 2006-2012 Kim Tore Jensen)
AC_CONFIG_SRCDIR([src/main.cpp])
AC_CONFIG_HEADERS([build.h])

# Checks for programs.
AC_PROG_CXX
AC_PROG_CC

# Checks for libraries.

# Checks for header files.
AC_CHECK_HEADERS([locale.h netdb.h stdlib.h sys/socket.h sys/time.h unistd.h])

# Checks for typedefs, structures, and compiler characteristics.
AC_HEADER_STDBOOL
AC_C_INLINE
AC_TYPE_SIZE_T

# Checks for library functions.
AC_CHECK_FUNCS([gettimeofday memset select setlocale socket])

AC_CONFIG_FILES([Makefile src/Makefile])
AC_OUTPUT
