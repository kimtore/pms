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


#include "input.h"
#include "config.h"
#include "pms.h"

extern Pms *		pms;


Input::Input()
{
	this->_mode = INPUT_NORMAL;
	pending = PEND_NONE;
	text.clear();
	searchterm.clear();
	cmdhistory.clear();
	searchhistory.clear();

	winclear();
}

Input::~Input()
{
}

/*
 * Store values when an extra window parameter is needed
 */
void			Input::winstore(pms_window * w)
{
	win = w;
	winparam = param;
	winpend = pending;
	pms->log(MSG_DEBUG, 0, "winstore: win=%p, winparam=%s, winpend=%d\n", win, winparam.c_str(), winpend);
}

/*
 * Delete window parameters
 */
void			Input::winclear()
{
	win = NULL;
	winparam = "";
	winpend = PEND_NONE;
}

/*
 * Restore window parameters
 */
bool			Input::winpop()
{
	if (win == NULL) return false;
	pms->log(MSG_DEBUG, 0, "winpop: setting param=%s, pending=%d\n", param.c_str(), pending);

	param = winparam;
	pending = winpend;

	return true;
}

/*
 * Go to next history item
 */
bool			Input::gonext()
{
	vector<string>::const_iterator	e;

	if (_mode == INPUT_JUMP || _mode == INPUT_FILTER)
		e = searchhistory.end();
	else if (_mode == INPUT_COMMAND)
		e = cmdhistory.end();
	else
		return false;

	if (historypos == e)
		return false;

	++historypos;

	if (historypos == e)
	{
		text.clear();
		return false;
	}
	text = *historypos;
	return true;
}

/*
 * Go to previous history item
 */
bool			Input::goprev()
{
	vector<string>::const_iterator	e;

	if (_mode == INPUT_JUMP || _mode == INPUT_FILTER)
		e = searchhistory.begin();
	else if (_mode == INPUT_COMMAND)
		e = cmdhistory.begin();
	else
		return false;

	if (historypos == e)
		return false;

	--historypos;
	text = *historypos;
	return true;
}

/*
 * Set input mode
 */
void			Input::mode(Input_mode m)
{
	pms->log(MSG_DEBUG, 0, "Entering input mode %d\n", m);

	switch(m)
	{
		case INPUT_COMMAND:
			text.clear();
			historypos = cmdhistory.end();
			break;

		case INPUT_JUMP:
		case INPUT_FILTER:
			text.clear();
			searchterm.clear();
			historypos = searchhistory.end();
			break;

		default:
			break;
	}

	this->_mode = m;
}

Input_mode		Input::mode()
{
	return this->_mode;
}

int			Input::get_keystroke()
{
	ch = getch();
	if (ch != -1)
	{
		return true;
	}

	return false;
}

pms_pending_keys	Input::dispatch()
{
	this->pending = PEND_NONE;
	this->param.clear();

	if (_mode == INPUT_NORMAL)
		return dispatch_normal();
	else if (_mode == INPUT_LIST)
		return dispatch_list();
	else
		return dispatch_text();

	return PEND_NONE;
}

pms_pending_keys	Input::dispatch_text()
{
	switch(ch)
	{
		case KEY_RESIZE:	/* Window was resized */
			pending = PEND_RESIZE;
			break;
		case 10:		/* Return */
		case 343:		/* Enter (keypad) */
			pending = PEND_TEXT_RETURN;
			break;
		case 21:		/* ^U */
			if (text.size() == 0)
			{
				text.clear();
				pending = PEND_TEXT_ESCAPE;
			}
			else
			{
				text.clear();
				pending = PEND_TEXT_UPDATED;
			}
			break;
		case 27:		/* Escape */
			pending = PEND_RETURN_ESCAPE;
			break;
		case 8:			/* ^H -- backspace */
		case 127:		/* ^? -- delete */
		case KEY_BACKSPACE:
			if (text.size() > 0)
			{
				text.erase(--string::iterator(text.end()));
				pending = PEND_TEXT_UPDATED;
			}
			else
			{
				pending = PEND_TEXT_ESCAPE;
			}
			break;
		case KEY_MOUSE:
			pending = PEND_NONE;
			break;
		default:
			if (ch < 32 || ch >= KEY_CODE_YES)
			{
				pending = pms->bindings->act(ch, &param);
				if (pending != PEND_NONE)
				{
					break;
				}
				pms->log(MSG_DEBUG, 0, "Key %3d '%c' pressed in text mode but not textual and not bound.\n", ch, ch);
			}
			text += ch;

			pending = PEND_TEXT_UPDATED;
	}

	return pending;
}

pms_pending_keys	Input::dispatch_list()
{
	if (ch == 10 || ch == 343)
	{
		pending = PEND_RETURN;
		return pending;
	}
	else if (ch == 27)
	{
		pending = PEND_RETURN_ESCAPE;
		return pending;
	}

	return dispatch_normal();
}

pms_pending_keys	Input::dispatch_normal()
{
	MEVENT			mouseevent;
	int			mousewinx, mousewiny;
	bool			mousecurwin = false;
	bool			mousetopbar = false;
	bool			mousestatusbar = false;
	bool			mousepositionreadout = false;
	bool			mousemodshift = false;
	bool			mousemodctrl = false;
	bool			mousemodalt = false;
	int			mouselistindex;

	if (ch == -1)
		return PEND_NONE;

	if (ch == KEY_RESIZE)
	{
		pending = PEND_RESIZE;
		return pending;
	}

	/* Mouse event */
	if (ch == KEY_MOUSE)
	{
		if (getmouse(&mouseevent) == ERR)
		{
			pms->log(MSG_DEBUG, 0, "error with getmouse()\n");
			ch = -1; // prevents weird results
			return PEND_NONE;
		}

		pms->log(MSG_DEBUG, 0, "mevent x:%d, y:%d, z:%d\n", mouseevent.x, mouseevent.y, mouseevent.z);

		if (mouseevent.bstate & BUTTON_SHIFT)
		{
			pms->log(MSG_DEBUG, 0, "shift is down\n");
			mousemodshift = true;
		}
		if (mouseevent.bstate & BUTTON_CTRL)
		{
			pms->log(MSG_DEBUG, 0, "ctrl is down\n");
			mousemodctrl = true;
		}
		if (mouseevent.bstate & BUTTON_ALT)
		{
			pms->log(MSG_DEBUG, 0, "alt is down\n");
			mousemodctrl = true;
		}

		mousewinx = mouseevent.x;
		mousewiny = mouseevent.y;

		if (pms->disp->actwin() && wenclose(pms->disp->actwin()->h(), mouseevent.y, mouseevent.x))
		{
			pms->log(MSG_DEBUG, 0, "mouse event in current window\n");
			mousecurwin = true;
			wmouse_trafo(pms->disp->actwin()->h(), &mousewiny, &mousewinx, false);

			//take window title and column titles away
			mousewiny -= 2;

			mouselistindex = pms->disp->actwin()->cursordrawstart() + mousewiny;
			if (!pms->disp->actwin()->plist() || mouselistindex >= pms->disp->actwin()->plist()->size())
			{
				//not a playlist or clicked off the end of the 
				//list
				mouselistindex = -1;
				pms->log(MSG_DEBUG, 0, "mouse event off the end of the list of songs or not a playlist\n");
			}
		}
		else if (wenclose(pms->disp->topbar->h(), mouseevent.y, mouseevent.x))
		{
			pms->log(MSG_DEBUG, 0, "mouse event in topbar\n");
			mousetopbar = true;
			wmouse_trafo(pms->disp->topbar->h(), &mousewiny, &mousewinx, false);
		}
		else if (wenclose(pms->disp->statusbar->h(), mouseevent.y, mouseevent.x))
		{
			pms->log(MSG_DEBUG, 0, "mouse event in statusbar\n");
			mousestatusbar = true;
			wmouse_trafo(pms->disp->statusbar->h(), &mousewiny, &mousewinx, false);
		}
		else if (wenclose(pms->disp->positionreadout->h(), mouseevent.y, mouseevent.x))
		{
			pms->log(MSG_DEBUG, 0, "mouse event in positionreadout\n");
			mousepositionreadout = true;
			wmouse_trafo(pms->disp->positionreadout->h(), &mousewiny, &mousewinx, false);
		}
		else
		{
			pms->log(MSG_DEBUG, 0, "mouse event doesn't seem to be enclosed in any of our windows\n");
			return PEND_NONE;
		}

		pms->log(MSG_DEBUG, 0, "mouse event at row %d, col %d of window\n", mousewiny, mousewinx);

		if (mouseevent.bstate & MOUSEWHEEL_DOWN)
		{
			pms->log(MSG_DEBUG, 0, "mousewheel down\n");
			if (mousetopbar)
			{
				if (mousemodctrl)
				{
					param = "-3";
					return PEND_VOLUME;
				}
				return PEND_NEXT;
			}
			if (mousecurwin)
			{
				if (mousewiny == -2) //heading bar
					return PEND_NEXTWIN;
				return PEND_SCROLL_DOWN;
			}
			return PEND_NONE;
		}
		else if (mouseevent.bstate & MOUSEWHEEL_UP)
		{
			pms->log(MSG_DEBUG, 0, "mousewheel up\n");
			if (mousetopbar)
			{
				if (mousemodctrl) {
					param = "+3";
					return PEND_VOLUME;
				}
				return PEND_PREV;
			}
			if (mousecurwin)
			{
				if (mousewiny == -2) //heading bar
					return PEND_PREVWIN;
				return PEND_SCROLL_UP;
			}
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON1_PRESSED)
		{
			pms->log(MSG_DEBUG, 0, "button 1 down\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON1_RELEASED)
		{
			pms->log(MSG_DEBUG, 0, "button 1 released\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON1_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 1 clicked\n");
			if (mousetopbar)
				return PEND_TOGGLEPLAY;
			if (mousecurwin)
			{
				if (mousewiny == -2) //heading bar
					return PEND_NEXTWIN;
				if (mouselistindex >= 0) //song
				{
					pms->disp->actwin()->plist()->setcursor(mouselistindex);
					if (mousemodctrl)
						pms->disp->actwin()->plist()->selectsong(pms->disp->actwin()->plist()->song(mouselistindex), !pms->disp->actwin()->plist()->song(mouselistindex)->selected);
					return PEND_NONE;
				}
			}
			if (mousestatusbar)
				return PEND_COMMANDMODE;
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON1_DOUBLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 1 doubleclicked\n");
			if (mousetopbar)
				return PEND_STOP;
			if (mousecurwin)
			{
				if (mousewiny == -2) //heading bar
					return PEND_PREVWIN;
				if (mouselistindex >= 0) //song
				{
					pms->disp->actwin()->plist()->setcursor(mouselistindex);
					return PEND_PLAY;
				}
			}
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON1_TRIPLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 1 tripleclicked\n");
			if (mousecurwin && mouselistindex >= 0)
			{
				pms->disp->actwin()->plist()->setcursor(mouselistindex);
				return PEND_ADD;
			}
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON2_PRESSED)
		{
			pms->log(MSG_DEBUG, 0, "button 2 down\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON2_RELEASED)
		{
			pms->log(MSG_DEBUG, 0, "button 2 released\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON2_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 2 clicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON2_DOUBLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 2 doubleclicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON2_TRIPLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 2 tripleclicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON3_PRESSED)
		{
			pms->log(MSG_DEBUG, 0, "button 3 down\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON3_RELEASED)
		{
			pms->log(MSG_DEBUG, 0, "button 3 released\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON3_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 3 clicked\n");
			if (mousecurwin)
			{
				if (mouselistindex >= 0) //song
				{
					pms->disp->actwin()->plist()->setcursor(mouselistindex);
					pms->disp->actwin()->plist()->selectsong(pms->disp->actwin()->plist()->song(mouselistindex), !pms->disp->actwin()->plist()->song(mouselistindex)->selected);
					return PEND_NONE;
				}
			}
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON3_DOUBLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 3 doubleclicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON3_TRIPLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 3 tripleclicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON4_PRESSED)
		{
			pms->log(MSG_DEBUG, 0, "button 4 down\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON4_RELEASED)
		{
			pms->log(MSG_DEBUG, 0, "button 4 released\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON4_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 4 clicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON4_DOUBLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 4 doubleclicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON4_TRIPLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 4 tripleclicked\n");
			return PEND_NONE;
		}
#if NCURSES_MOUSE_VERSION > 1
		else if (mouseevent.bstate & BUTTON5_PRESSED)
		{
			pms->log(MSG_DEBUG, 0, "button 5 down\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON5_RELEASED)
		{
			pms->log(MSG_DEBUG, 0, "button 5 released\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON5_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 5 clicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON5_DOUBLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 5 doubleclicked\n");
			return PEND_NONE;
		}
		else if (mouseevent.bstate & BUTTON5_TRIPLE_CLICKED)
		{
			pms->log(MSG_DEBUG, 0, "button 5 tripleclicked\n");
			return PEND_NONE;
		}
#endif
		else if (mouseevent.bstate & REPORT_MOUSE_POSITION)
		{
			pms->log(MSG_DEBUG, 0, "mouse position -- what does this do?\n");
			return PEND_NONE;
		}
		else
		{
			pms->log(MSG_DEBUG, 0, "mevent state (%d) unknown\n", mouseevent.bstate);
			return PEND_NONE;
		}
	}

	/* Key pressed */
	pending = pms->bindings->act(ch, &param);

	if (pending == PEND_NONE)
	{
		pms->log(MSG_STATUS, STERR, _("Key is not bound."));
		pms->log(MSG_DEBUG, 0, "Key %3d '%c' pressed but not bound.\n", ch, ch);
	}

	return pending;
}

/*
 * Saves the current text in history
 */
void		Input::savehistory()
{
	switch(_mode)
	{
		case INPUT_JUMP:
		case INPUT_FILTER:
			searchhistory.push_back(text);
			historypos = searchhistory.end();
			break;

		case INPUT_COMMAND:
			cmdhistory.push_back(text);
			historypos = cmdhistory.end();
			break;

		default:
			break;
	}
}

/*
 * Run a command
 */
bool		Input::run(string s, Message & err)
{
	int		pos;

	if (s.size() == 0)
		return true;

	pos = s.find_first_of(" ");
	if (pos > 0)
	{
		param = s.substr(pos + 1);
		s = s.substr(0, pos);
		pms->log(MSG_DEBUG, 0, "Running command '%s' with param '%s'\n", s.c_str(), param.c_str());
	}
	else
		pms->log(MSG_DEBUG, 0, "Running command '%s' without parameters\n", s.c_str());

	pending = pms->commands->act(s);

	if (pending == PEND_NONE)
	{
		err.code = CERR_UNKNOWN_COMMAND;
		err.str = _("unknown command");
		err.str += " '" + s + "'";
		return false;
	}

	return true;
}
