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

#ifndef _SET_PARAMETERS_H_
#define _SET_PARAMETERS_H_

#include "options.h"

using namespace std;


/**
 * The option_operator_t enum denotes how to operate on a variable.
 *
 * - OPTION_OPERATOR_SET: set the variable to a non-boolean value.
 * - OPTION_OPERATOR_INVERT: invert the boolean variable.
 * - OPTION_OPERATOR_SET_TRUE: set the boolean variable to true.
 * - OPTION_OPERATOR_SET_FALSE: set the boolean variable to false.
 * - OPTION_OPERATOR_SET_QUERY: use this operator in order to query for the
 *   existing value of the option.
 */
typedef enum
{
	OPTION_OPERATOR_SET,
	OPTION_OPERATOR_INVERT,
	OPTION_OPERATOR_SET_TRUE,
	OPTION_OPERATOR_SET_FALSE,
	OPTION_OPERATOR_QUERY
}
option_operator_t;


/**
 * The SetParameters class parses a "option=value" argument into normalized
 * components, and may set or retrieve the existing value as a string value.
 */
class SetParameters
{
private:
	/**
	 * Parse a "option=value" string into normalized components.
	 */
	void			parse(string input);

	/**
	 * Return `value' as a boolean.
	 */
	bool			make_bool();

	/**
	 * Return `value' as a long integer.
	 */
	bool			make_long();

	/**
	 * Converted values from make_bool() and make_long().
	 */
	bool			value_bool;
	long			value_long;

	/**
	 * Normalized components.
	 */
	string			name;
	string			value;
	option_operator_t	operator_;

	/**
	 * Matching Option instance.
	 */
	Option *		option;


public:
	/**
	 * Instantiate the object and parse the input string.
	 */
	SetParameters(string input);

	/**
	 * Copy the option value to the Option object, using its
	 * native data type.
	 */
	bool			commit();

	/**
	 * Return true if this option is real, and can be looked up in the
	 * Options class.
	 */
	bool			exists();

	/**
	 * Return true if this is an option query (i.e., do not store the
	 * value, only retrieve it), and the option exists.
	 */
	bool			is_valid_query();

	/**
	 * Return true if the operator is valid for this data type.
	 */
	bool			has_valid_operator();

	/**
	 * Copy the option value into the value string buffer.
	 *
	 * Returns true on success, false on failure.
	 */
	bool			load();

	/**
	 * Produce a string representation of this option value.
	 *
	 * In case of a non-boolean option, returns the "name=value" pair as a string.
	 *
	 * In case of a boolean option, returns "name" or "noname" if the value is true
	 * or false, respectively.
	 */
	string			repr();
};

#endif /* _SET_PARAMETERS_H_ */
