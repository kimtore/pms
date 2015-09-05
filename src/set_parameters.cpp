/* vi:set ts=8 sts=8 sw=8 noet:
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
 */

#include <errno.h>

#include "set_parameters.h"
#include "pms.h"

using namespace std;

extern Pms * pms;

SetParameters::SetParameters(const string input)
{
	name.clear();
	value.clear();
	operator_ = OPTION_OPERATOR_SET_TRUE;
	option = NULL;

	parse(input);
}

void
SetParameters::parse(const string input)
{
	string::const_iterator	it;
	bool			processed_sign = false;

	it = input.begin();

	/* Parse the string into name, value, and sign components. */
	while (it != input.end()) {
		if (!processed_sign && (*it == ':' || *it == '=')) {
			operator_ = OPTION_OPERATOR_SET;
			processed_sign = true;
		} else if (!processed_sign) {
			name += *it;
		} else if (processed_sign) {
			value += *it;
		}
		++it;
	}

	/* At least a name is needed to continue parsing. */
	if (!name.size()) {
		return;
	}

	/* Check if the option exists, and stop parsing if it does. */
	if ((option = pms->options->lookup_option(name.c_str())) != NULL) {
		return;
	}

	/* Option might be prefixed by "inv" or "no",
	 * or suffixed by "!" or "?". Let's handle those cases. */
	if (name.substr(0, 2) == "no") {
		operator_ = OPTION_OPERATOR_SET_FALSE;
		name = name.substr(2);
	} else if (name.substr(0, 3) == "inv") {
		operator_ = OPTION_OPERATOR_INVERT;
		name = name.substr(3);
	} else if (name.substr(name.size() - 1, 1) == "!") {
		operator_ = OPTION_OPERATOR_INVERT;
		name = name.substr(0, name.size() - 1);
	} else if (name.substr(name.size() - 1, 1) == "?") {
		operator_ = OPTION_OPERATOR_QUERY;
		name = name.substr(0, name.size() - 1);
	} else {
		return;
	}

	option = pms->options->lookup_option(name.c_str());
}

bool
SetParameters::exists()
{
	return (option != NULL);
}

bool
SetParameters::is_valid_query()
{
	return (exists() && operator_ == OPTION_OPERATOR_QUERY);
}

bool
SetParameters::has_valid_operator()
{
	/* Cannot perform boolean operations on non-boolean options. */
	if (option->type != OPT_BOOL && operator_ != OPTION_OPERATOR_SET) {
		return false;
	}

	/* Vice versa. */
	if (option->type == OPT_BOOL && operator_ == OPTION_OPERATOR_SET) {
		return false;
	}

	return true;
}

bool
SetParameters::load()
{
	if (!exists()) {
		return false;
	}

	if (option->type == OPT_LONG) {
		value = Pms::tostring(*(option->as_long_ptr()));
	} else if (option->type == OPT_STRING) {
		value = *(option->as_string_ptr());
	} else {
		return false;
	}

	return true;
}

string
SetParameters::repr()
{
	if (!exists()) {
		return "";
	}

	if (option->type == OPT_BOOL) {
		if (*(option->as_bool_ptr())) {
			return name;
		}
		return "no" + name;
	}

	if (!load()) {
		return "";
	}

	return name + "=" + value;
}

bool
SetParameters::commit()
{
	if (!exists() || !has_valid_operator()) {
		return false;
	}

	switch(option->type) {
		case OPT_BOOL:
			if (!make_bool()) {
				return false;
			}
			return option->set(&value_bool);
		case OPT_LONG:
			if (!make_long()) {
				return false;
			}
			return option->set(&value_long);
		case OPT_STRING:
			return option->set(&value);
		default:
			abort();
	}
}

bool
SetParameters::make_bool()
{
	switch(operator_) {
		case OPTION_OPERATOR_INVERT:
			value_bool = !(*(static_cast<bool *>(option->pointer)));
			break;
		case OPTION_OPERATOR_SET_TRUE:
			value_bool = true;
			break;
		case OPTION_OPERATOR_SET_FALSE:
			value_bool = false;
			break;
		default:
			abort();
	}

	return true;
}

bool
SetParameters::make_long()
{
	long temp_value;
	const char * nptr = value.c_str();
	char * endptr;

	errno = 0;
	temp_value = strtol(nptr, &endptr, 0);

	/* No conversion performed */
	if (nptr == endptr) {
		return false;
	}

	/* An error occurred, also includes underflow and overflow errors */
	if (errno != 0) {
		return false;
	}

	value_long = temp_value;
	return true;
}
