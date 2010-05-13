/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2010  Kim Tore Jensen
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


#ifndef _INPUT_H_
#define _INPUT_H_

#include "mycurses.h"
#include "display.h"
#include "types.h"
#include "message.h"

#if NCURSES_MOUSE_VERSION > 1
#define MOUSEWHEEL_DOWN	BUTTON4_PRESSED | BUTTON4_CLICKED | BUTTON4_DOUBLE_CLICKED | BUTTON4_TRIPLE_CLICKED
#define MOUSEWHEEL_UP	BUTTON5_PRESSED | BUTTON5_CLICKED | BUTTON5_DOUBLE_CLICKED | BUTTON5_TRIPLE_CLICKED
#else
//mousewheel isn't supposed to be supported, but this works for me (tremby). 
//however, it does break the button 2 pressed event (pressing the wheel button)
#define MOUSEWHEEL_DOWN	(BUTTON2_PRESSED | REPORT_MOUSE_POSITION)
#define MOUSEWHEEL_UP	BUTTON4_PRESSED
#endif

typedef enum Input_mode
{
	INPUT_NORMAL,
	INPUT_JUMP,
	INPUT_FILTER,
	INPUT_COMMAND,
	INPUT_LIST

} Input_mode;


class Input
{
private:
	Input_mode			_mode;
	pms_pending_keys		pending;
	wchar_t				ch;
	vector<string>			cmdhistory;
	vector<string>			searchhistory;
	vector<string>::iterator	historypos;

	pms_pending_keys		dispatch_normal();
	pms_pending_keys		dispatch_list();
	pms_pending_keys		dispatch_text();

public:
	Input();
	~Input();

	string			param;
	string			text;
	string			searchterm;

	/* Storing values when a window parameter is needed */
	string			winparam;
	pms_pending_keys	winpend;
	pms_window *		win;

	void			winstore(pms_window *);
	void			winclear();
	bool			winpop();

	bool			gonext();	// Go to next history item
	bool			goprev();	// Go to previous history item

	void			mode(Input_mode);
	Input_mode		mode();
	int			get_keystroke();
	pms_pending_keys	dispatch();
	pms_pending_keys	getpending() { return pending; };
	void			savehistory();
	bool			run(string, Message &);

};


#endif /* _INPUT_H_ */
