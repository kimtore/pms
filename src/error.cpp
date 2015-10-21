/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2015  Kim Tore Jensen
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#include <cstdarg>
#include <cstdio>
#include <cstring>

#include "error.h"

static char pms_errstr[1024];

void
pms_error(const char * format, ...)
{
	va_list ap;
        char errstr[1024];
	char buffer[1024];

	va_start(ap, format);
	vsnprintf(buffer, sizeof(buffer) - 1, format, ap);
	va_end(ap);

        if (strlen(pms_errstr)) {
                snprintf(errstr, sizeof(errstr) - 1, "%s: %s", buffer, pms_errstr);
                strcpy(buffer, errstr);
        }

        strncpy(pms_errstr, buffer, sizeof(pms_errstr) - 1);
}

void
pms_error_clear()
{
        pms_errstr[0] = '\0';
}

const char *
pms_error_string()
{
        return pms_errstr;
}
