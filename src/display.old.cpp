/*
 * Draw the position readout
 */
void			pms_win_positionreadout::draw()
{
	pms_window *	win = pms->disp->actwin();
	string		text;
	char		chararray[4];

	if (!win)
		text = "";
	else if (win->size() <= win->bheight() - 1)
		text = "All";
	else if (win->cursordrawstart() == 0)
		text = "Top";
	else if (win->cursordrawstart() == win->size() - (win->bheight() - 1))
		text = "Bot";
	else
	{
		sprintf(chararray, "%2d", 100 * win->cursordrawstart() / (win->size() - (win->bheight() - 1)));
		text = chararray;
		text += "%%";
	}

	/* Clear window */
	clear(false, 0);

	/* Draw string */
	colprint(this, 0, 0, pms->options->colors->position, "%s", text.c_str());

	return;
}

/*
 *
 * Topbar window class
 *
 */
pms_win_topbar::pms_win_topbar(Control * c)
{
	comm = c;
}

/*
 * Draws up-to-date info about current song.
 */
void			pms_win_topbar::draw()
{
	unsigned int	y, x, reallen;
	int		drawx, drawlen, i, progress;
	string		t;
	Songlist *	list;
	Song *		song;

	/* No-go */
	if (pms->options->topbar_lines.size() == 0 || !pms->options->topbarvisible)
		return;

	/* Clear window */
	clear(false, 0);
	wantdraw = false;

	/* Draw info from topbar class */
	song = pms->cursong();
	for (y = 0; y < pms->options->topbar_lines.size(); y++)
	{
		x = 0;
		while (true)
		{
			t = pms->formatter->format(song, pms->options->topbar_lines[y]->strings[x], reallen, &(pms->options->colors->topbar.fields));
			if (reallen != 0 && t.size() != 0)
			{
				drawlen = static_cast<int>(reallen);

				if (x == 0)
					drawx = 0;
				else if (x == 1)
					drawx = (bwidth() / 2) - (drawlen / 2);
				else if (x == 2)
					drawx = (bwidth() - drawlen);

				colprint(this, y, drawx, pms->options->colors->topbar.standard, t.c_str());
			}

			if (x == 0)
				x = 2;
			else if (x == 2)
				x = 1;
			else
				break;
		}
	}

	drawborders();

	return;
}


/*
 * Returns the height of the topbar window including borders and space if 
 * applicable
 */
int			pms_win_topbar::height()
{
	return pms->options->topbar_lines.size() + (pms->options->topbarborders ? 2 : 0);
}


/*
 *
 * End of topbar window class.
 *
 * Windowlist window class
 *
 */
pms_win_windowlist::pms_win_windowlist(Display * ndisp, vector<pms_win_playlist *> * wl)
{
	unsigned int		i;

	column.push_back(new pms_column("Songs", EINVALID, 7));
	column.push_back(new pms_column("Window title", EINVALID, 0));

	mydisp = ndisp;
	originwin = dynamic_cast<pms_win_playlist *>(ndisp->actwin());
	clastwin = pms->disp->lastwin;
	wlist = wl;

	for (i = 0; i < wlist->size(); i++)
	{
		if (static_cast<void *>((*wlist)[i]) == static_cast<void *>(mydisp->actwin()))
		{
			cursor = i;
			return;
		}
	}
}

/*
 * Switch cursor between last used windows
 */
void			pms_win_windowlist::switchlastwin()
{
	pms_window *		w;
	pms_window *		v;
	unsigned int		i;

	w = current();
	if (w == clastwin)
		w = dynamic_cast<pms_window *>(originwin);
	else
		w = clastwin;

	for (i = 0; i < size(); i++)
	{
		v = dynamic_cast<pms_window *>((*wlist)[i]);
		if (w == v)
		{
			setcursor(i);
			return;
		}
	}
}

/*
 * Return window under cursor
 */
pms_window *		pms_win_windowlist::current()
{
	if (cursor < 0)
		cursor = 0;
	if (cursor >= static_cast<int>(wlist->size()))
		cursor = static_cast<int>(wlist->size()) - 1;

	if (cursor < 0)
		return NULL;

	return dynamic_cast<pms_window *>((*wlist)[cursor]);
}

/*
 * Return last active window
 */
pms_window *		pms_win_windowlist::lastwin()
{
	return (mydisp ? mydisp->actwin() : NULL);
}

/*
 * Draw a list of all windows
 */
void			pms_win_windowlist::draw()
{
	unsigned int		i, j;
	unsigned int		min, max;
	int			songcount;
	unsigned int		counter = 0;
	color *			hilight;
	string			t;
	pms_win_playlist *	w;
	pms_win_playlist *	activewin;

	/* Clear window first */
	clear(false, 0);
	wantdraw = false;
	if (!wlist) return;

	min	= cursordrawstart();
	max	= min + (unsigned int)(bheight() - 1);
	if (max > wlist->size())
		max = wlist->size();

	activewin = dynamic_cast<pms_win_playlist *>(mydisp->playingwin());

	/* Traverse window list and draw */
	for (i = 0; i < wlist->size(); i++)
	{
		w = (*wlist)[i];
		if (!w) continue;
		if ((void *)w == (void *)this) continue;

		++counter;
		if (i == static_cast<unsigned int>(cursor))
			hilight = pms->options->colors->cursor;
		else if (w == selected)
			hilight = pms->options->colors->selection;
		else if (w == originwin)
			hilight = pms->options->colors->lastlist;
		else if (w == activewin)
			hilight = pms->options->colors->playinglist;
		else
			hilight = NULL;

		if (w->plist())
			songcount = w->plist()->size();
		else
			songcount = -1;

		if (hilight)
		{
			wattron(handle, hilight->pair());
			mvwhline(handle, counter + border[0], 0, ' ', COLS);
			wattroff(handle, hilight->pair());
		}

		/* Draw song count, if any */
		if (songcount >= 0)
		{
			t = Pms::tostring(songcount);
			colprint(this, counter, 0, (hilight ? hilight : pms->options->colors->fields.num), "%s", t.c_str());
		}

		/* Draw window title */
		t = w->fulltitle();
		colprint(this, counter, column[0]->len() + 1, hilight, "%s", t.c_str());
	}

	/* Draw captions and column borders */
	//TODO: make into an option
	j = 0;
	for (i = 0; i < column.size(); i++)
	{
		colprint(this, 0, (i == 0 ? j : j + 1),
			pms->options->colors->headers,
			"%s", column[i]->title.c_str());
		if (i > 0 && pms->options->columnborders)
		{
			wattron(handle, pms->options->colors->border->pair());
			mvwvline(handle, border[0], j, ACS_VLINE, bheight());
			wattroff(handle, pms->options->colors->border->pair());
		}
		j += column[i]->len();
	}

	drawborders();
}

/*
 * Return offset in song list of the topmost song visible
 */
unsigned int		pms_window::cursordrawstart()
{
	static float		f;
	static float		ht;
	int			i;
	unsigned int		cursordrawstart;
	int			sotemp;
	int			scrolloffsetmax;

	/* By default, always start at top */
	cursordrawstart = 0;

	/* Do nothing if window is empty */
	if (size() == 0)
	{}

	/* Cursors position on screen changes relative to cursor position in list */
	else if (pms->options->scroll_mode == SCROLL_RELATIVE)
	{
		ht	= bheight() - 2;

		if (ht >= size())
			cursordrawstart = 0;
		else
		{
			f = ((scursor() / float(size() - 1)) * ht);
			cursordrawstart = (static_cast<float>(scursor()) - round(f));
		}
	}

	/* Cursor is always centered, except when nearing top or bottom of the list */
	else if (pms->options->scroll_mode == SCROLL_CENTERED)
	{
		if (size() > static_cast<unsigned int>(bheight() - 1))
		{
			i = scursor() - (bheight() / 2) + 1;

			if (i < 0)
				cursordrawstart = 0;
			else if (i > static_cast<int>(size()) - bheight())
				cursordrawstart = (size() - static_cast<unsigned int>(bheight()) + 1);
			else
				cursordrawstart = static_cast<unsigned int>(i);
		}
	}

	/* Window is scrolled when the cursor is about to go off the edge */
	else if (pms->options->scroll_mode == SCROLL_NORMAL)
	{
		//note bheight() includes the column headings!

		//if scrolloff is set to half the height or more drop it 
		//temporarily
		sotemp = pms->options->scrolloff;
		if (sotemp * 2 >= bheight() - 1)
			sotemp = (bheight() - 1 - 1) / 2;

		//get rid of any empty space at the bottom which shouldn't be there
		while (scrolloffset > 0 && scrolloffset + bheight() - 1 > size())
			scrolloffset--;

		//is the cursor too high?
		i = scrolloffset + sotemp - scursor();
		if (i > 0)
		{
			scrolloffset -= i;
			if (scrolloffset < 0)
				scrolloffset = 0;
		}
		else
		{
			//is the cursor too low?
			i = scursor() - (scrolloffset + bheight() - 2 - sotemp);
			if (i > 0)
			{
				scrolloffset += i;

				scrolloffsetmax = static_cast<int>(size()) - (bheight() - 1);
				if (scrolloffsetmax < 0)
					scrolloffsetmax = 0;

				if (scrolloffset > scrolloffsetmax)
					scrolloffset = scrolloffsetmax;
			}
		}

		cursordrawstart = static_cast<unsigned int>(scrolloffset);
	}

	scrolloffset = static_cast<int>(cursordrawstart);
	return cursordrawstart;
}

/*
 * Set absolute cursor position
 */
void		pms_window::setcursor(int absolute)
{
	cursor = absolute;
	if (cursor < 0)
	{
		cursor = 0;
	}
	else if (cursor >= (int)size())
	{
		cursor = (int)(size() - 1);
	}

	wantdraw = true;
}

/*
 * Scroll window
 */
void		pms_window::scrollwin(int offset)
{
	int	i;
	int	sotemp;

	if (pms->options->scroll_mode != SCROLL_NORMAL)
	{
		movecursor(offset);
		return;
	}

	if (size() <= bheight() - 1)
	{
		return;
	}

	//if scrolloff is set to half the height or more drop it temporarily
	sotemp = pms->options->scrolloff;
	if (sotemp * 2 >= bheight() - 1)
		sotemp = (bheight() - 1 - 1) / 2;

	if (offset == 0)
		return;
	else if (offset < 0) //up
	{
		i = -cursordrawstart();
		if (offset < i)
			offset = i;
		if (offset == 0) {
			return;
		}

		//cursor too low?
		i = cursordrawstart() + (bheight() - 1) - 1 - sotemp + offset - scursor();
		if (i < 0)
			movecursor(i);
	}
	else //down
	{
		i = size() - (bheight() - 1) - cursordrawstart();
		if (offset > i)
			offset = i;
		if (offset == 0) {
			return;
		}

		//cursor too high?
		i = cursordrawstart() + sotemp + offset - scursor();
		if (i > 0)
			movecursor(i);
	}

	scrolloffset += offset;

	wantdraw = true;
}


/*
 *
 * End of window class.
 * Bindings window class.
 *
 */
pms_win_bindings::pms_win_bindings()
{
	column.push_back(new pms_column(_("Key"), EINVALID, 14));
	column.push_back(new pms_column(_("Command"), EINVALID, 30));
	column.push_back(new pms_column(_("Description"), EINVALID, 0));

	pms->bindings->list(&key, &command, &desc);
}

/*
 * Draw keypms->bindings
 */
void		pms_win_bindings::draw()
{
	unsigned int		counter = 0;
	unsigned int		i, j;
	unsigned int		min, max;
	color *			hilight;

	/* Clear window first */
	clear(false, 0);
	wantdraw = false;

	min	= cursordrawstart();
	max	= min + static_cast<unsigned int>(bheight() - 1);

	if (max >= key.size())
		max = key.size();

	/* Traverse pms->bindings and draw */
	for (i = min; i < max; i++)
	{
		++counter;
		if (i == cursor)
			hilight = pms->options->colors->cursor;
		else
			hilight = NULL;

		if (hilight)
		{
			wattron(handle, hilight->pair());
			mvwhline(handle, counter + border[0], 0, ' ', COLS);
			wattroff(handle, hilight->pair());
		}
		else
		{
			mvwhline(handle, counter + border[0], 0, ' ', COLS);
		}

		colprint(this, counter, 0, hilight, "%s", key[i].c_str());
		colprint(this, counter, column[0]->len() + 1, hilight, "%s", command[i].c_str());
		colprint(this, counter, column[0]->len() + column[1]->len() + 1, hilight, "%s", desc[i].c_str());
	}

	/* Draw captions and column borders */
	j = 0;
	for (i = 0; i < column.size(); i++)
	{
		colprint(this, 0, (i == 0 ? j : j + 1),
			pms->options->colors->headers,
			"%s", column[i]->title.c_str());
		if (i > 0 && pms->options->columnborders)
		{
			wattron(handle, pms->options->colors->border->pair());
			mvwvline(handle, border[0], j, ACS_VLINE, bheight());
			wattroff(handle, pms->options->colors->border->pair());
		}
		j += column[i]->len();
	}

	drawborders();
}


/*
 * End of pms->bindings window class.
 * Playlist window class
 *
 */
pms_win_playlist::pms_win_playlist()
{
	list = NULL;
}

/*
 * Returns the window title, with playlist name if any
 */
string		pms_win_playlist::fulltitle()
{
	string		t;

	if (title.size() || !list)
	{
		t = title;
	}
	else
	{
		t = "Playlist: ";
		if (!list->filename.size())
			t += "[Untitled]";
		else
			t += list->filename;
	}

	if (list != NULL && list->filtercount() > 0)
	{
		t += " <" + Pms::tostring(static_cast<size_t>(list->filtercount()));
		if (list->filtercount() > 1)
			t += _(" filters enabled");
		else
			t += _(" filter enabled");
		t += ">";
	}

	return t;
}

/*
 * Moves cursor position to current song
 */
bool		pms_win_playlist::gotocurrent()
{
	if (!plist()) return false;

	return plist()->gotocurrent();
}

/*
 * Move cursor in either direction
 */
void		pms_win_playlist::movecursor(int offset)
{
	if (!list) return;
	list->movecursor(offset);
	wantdraw = true;
}

/*
 * Set absolute cursor position
 */
void		pms_win_playlist::setcursor(int absolute)
{
	if (!list) return;
	list->setcursor(absolute);
	wantdraw = true;
}

/*
 * Sets current playlist
 */
void		pms_win_playlist::setplist(Songlist * l)
{
	list = l;
	set_column_size();
}

/*
 * Get position of a specific song but don't jump to it
 */
int		pms_win_playlist::posof_jump(string term, int start, bool reverse)
{
	int		i;

	if (!list || list->size() == 0) return MATCH_FAILED;
	i = start - 1;
	if (i < 0) i = list->end();
	if (reverse)
		i = list->match(term, start, i, MATCH_ALL | MATCH_REVERSE);
	else
		i = list->match(term, start, i, MATCH_ALL);

	return i;
}

/*
 * Jump to a specific song
 */
bool		pms_win_playlist::jumpto(string term, int start, bool reverse)
{
	int i = posof_jump(term, start, reverse);

	if (i == -1) return false;
	
	list->setcursor(i);
	wantdraw = true;

	return true;
}

/*
 * Draws the current songlist.
 */
void		pms_win_playlist::draw()
{
	unsigned int		pair;
	unsigned int		counter = 0;
	unsigned int		i, j, winlen;
	unsigned int		min;
	unsigned int		max;
	int			ii;
	Song			*song;
	string			t;
	color *			hilight;
	color *			c;

	/* Clear window first */
	clear(false, 0);
	wantdraw = false;
	if (!list) return;

	/* Define which items to draw */
	min	= cursordrawstart();
	max	= min + (unsigned int)(bheight() - 1);
	if (max > list->size())
		max = list->size();

	/* Traverse song list and draw lines */
	for (i = min; i < max; i++)
	{
		++counter;
		hilight = NULL;

		song = list->song(i);
		if (i == list->cursor())
		{
			hilight = pms->options->colors->cursor;
		}
		else if (song->selected)
		{
			hilight = pms->options->colors->selection;
		}
		else if (pms->cursong()) {
                        if ((list->role == LIST_ROLE_MAIN && pms->cursong()->id == song->id) || (list->role != LIST_ROLE_MAIN && song->file == pms->cursong()->file)) {
				hilight = pms->options->colors->current;
			}
		}

		winlen = 0;
		for (j = 0; j < column.size(); j++)
		{
			pair = 0;

                        /* Draw highlight line */
			if (hilight) wattron(handle, hilight->pair());
			mvwhline(handle, counter + border[0], winlen, ' ', column[j]->len() + 1);
			if (hilight) wattroff(handle, hilight->pair());
			
			c = pms->formatter->getcolor(column[j]->type, &(pms->options->colors->fields));
			if (c)
			{
				t = pms->formatter->format(song, column[j]->type);
				colprint(this, counter, (j == 0 ? winlen : winlen + 1),
					(hilight ? hilight : c),
					"%s", t.c_str());

			}

			winlen += column[j]->len();
		}

		hilight = pms->options->colors->standard;
	}

	/* Draw captions and column borders */
	j = 0;
	for (i = 0; i < column.size(); i++)
	{
		colprint(this, 0, (i == 0 ? j : j + 1),
			pms->options->colors->headers,
			"%s", column[i]->title.c_str());
		if (i > 0 && pms->options->columnborders)
		{
			wattron(handle, pms->options->colors->border->pair());
			mvwvline(handle, border[0], j, ACS_VLINE, bheight());
			wattroff(handle, pms->options->colors->border->pair());
		}
		j += column[i]->len();
	}

	drawborders();
}


/*
 *
 * End of playlist window class.
 *
 * Window base class
 *
 */
pms_window::pms_window()
{
	handle = NULL;

	x = -1;
	y = -1;
	width = -1;
	height = -1;
	wantdraw = false;
	cursor = 0;
	scrolloffset = 0;

	memset(&border, 0, sizeof(border));
}

pms_window::~pms_window()
{
	delwin(this->handle);
}

/*
 * Set new window coordinates.
 */
bool		pms_window::resize(int nx, int ny, int nwidth, int nheight)
{
	if (nx < 0 || ny < 0 || nwidth <= 0 || nheight <= 0)
	{
		return false;
	}

	if (handle != NULL)
		delwin(this->handle);

	handle = newwin(nheight, nwidth, ny, nx);

	if (handle == NULL)
	{
		pms->log(MSG_DEBUG, 0, "resize: window creation FAILED: (%d, %d, %d, %d), exiting\n", nx, ny, nwidth, nheight);
		return false;
	}

	x = nx;
	y = ny;
	width = nwidth;
	height = nheight;

	wantdraw = true;

	return true;
}

/*
 * Draw whitespace over entire window, except borders.
 */
void		pms_window::clear(bool clearborders, color * c)
{
	int		y;

	if (c != NULL)
		wattron(handle, c->pair());
	
	for (y = border[0]; y < (clearborders ? height : height - border[2]); y++)
	{
		if (clearborders == true)
			mvwhline(handle, y, 0, ' ', width);
		else
			mvwhline(handle, y, border[3], ' ', width - border[1]);
	}

	if (c != NULL)
		wattroff(handle, c->pair());
}

/*
 * Set new window title
 */
void		pms_window::settitle(string ntitle)
{
	title = ntitle;
	wantdraw = true;
}

/*
 * Draw window title at top position
 */
void		pms_window::drawtitle()
{
	int		left, right;
	string		t;

	if (!handle) return;

	t = fulltitle();
	if (t.size() == 0)
		return;

	left = centered(t);
	right = left + t.size();

	wattron(handle, pms->options->colors->border->pair());
	mvwaddch(handle, 0, left - 2, ACS_RTEE);
	mvwaddch(handle, 0, right + 1, ACS_LTEE);
	wattroff(handle, pms->options->colors->border->pair());

	wattron(handle, pms->options->colors->title->pair());
	mvwprintw(handle, 0, left - 1, " %s ", t.c_str());
	wattroff(handle, pms->options->colors->title->pair());
}

/*
 * Set new borders, and draw them.
 */
void		pms_window::setborders(bool settop, bool setright, bool setbottom, bool setleft)
{
	if (!handle) return;

	wattron(handle, pms->options->colors->border->pair());

	border[0] = (settop ? 1 : 0);
	border[1] = (setright ? 1 : 0);
	border[2] = (setbottom ? 1 : 0);
	border[3] = (setleft ? 1 : 0);

	if (settop)
		mvwhline(handle, 0, 0, ACS_HLINE, COLS);
	if (setright)
		mvwvline(handle, 0, width-1, ACS_VLINE, LINES);
	if (setbottom)
		mvwhline(handle, height-1, 0, ACS_HLINE, COLS);
	if (setleft)
		mvwvline(handle, 0, 0, ACS_VLINE, LINES);

	if (settop && setright)		mvwaddch(handle, 0, width-1, ACS_URCORNER);
	if (settop && setleft)		mvwaddch(handle, 0, 0, ACS_ULCORNER);
	if (setbottom && setright)	mvwaddch(handle, height-1, width-1, ACS_LRCORNER);
	if (setbottom && setleft)	mvwaddch(handle, height-1, 0, ACS_LLCORNER);

	wattroff(handle, pms->options->colors->border->pair());

	if (border[0])		drawtitle();
}

/*
 * Draw window borders.
 */
void		pms_window::drawborders()
{
	setborders(border[0], border[1], border[2], border[3]);
}








/*
 * Return window containing the active playlist
 */
pms_window *		Display::playingwin()
{
	return findwlist(comm->activelist());
}



/*
 * Finds window with specified playlist
 */
pms_window *		Display::findwlist(Songlist * target)
{
	vector<pms_window *>::iterator	i;

	i = windows.begin();
	while (i != windows.end())
	{
		if ((*i)->plist() == target)
			return *i;
		++i;
	}

	return NULL;
}

/*
 * Returns a pointer to the next window
 */
pms_window *		Display::nextwindow()
{
	vector<pms_window *>::iterator	i;

	if (windows.size() == 0)
		return NULL;

	i = windows.begin();
	while (i != windows.end())
	{
		if (*i == curwin)
		{
			++i;
			if (i == windows.end())
			{
				i = windows.begin();
			}
			return *i;
		}
		++i;
	}

	return NULL;
}

/*
 * Returns a pointer to the previous window
 */
pms_window *		Display::prevwindow()
{
	vector<pms_window *>::iterator	i;

	if (windows.size() == 0)
		return NULL;

	i = windows.begin();
	while (i != windows.end())
	{
		if (*i == curwin)
		{
			if (i == windows.begin())
				i = windows.end();

			--i;
			return *i;
		}
		++i;
	}

	return NULL;
}

/*
 * Activates a window for input
 */
bool			Display::activate(pms_window * w)
{
	vector<pms_window *>::iterator	i;

	i = windows.begin();
	while (i != windows.end())
	{
		if (*i == w)
		{
			if (curwin && curwin->type() == WIN_ROLE_PLAYLIST)
			{
				lastwin = curwin;
				pms->log(MSG_DEBUG, 0, "Activate: setting lastwin=%p with list %p.\n", lastwin, lastwin->plist());
			}
			if (curwin && curwin != w)
			{
				switch(curwin->type())
				{
					case WIN_ROLE_WINDOWLIST:
					case WIN_ROLE_BINDLIST:
						delete_window(curwin);
						break;
					default:
						break;
				}
			}
			curwin = w;
			curwin->wantdraw = true;
			return true;
		}
		++i;
	}

	return false;
}

/*
 * Creates the key pms->bindings window
 */
pms_win_bindings *	Display::create_bindlist()
{
	pms_win_bindings *	w;

	w = new pms_win_bindings();
	if (w)
	{
		if (pms->options->topbar_lines.size() == 0 || !pms->options->topbarvisible)
			w->resize(0, 0, COLS, LINES - 1);
		else
			w->resize(0, pms->disp->topbar->height(), COLS, LINES - pms->disp->topbar->height() - 1);
		w->setborders(true, false, false, false);
		w->settitle("Key bindings");
		windows.push_back(w);
		if (curwin == NULL)
			curwin = w;
	}

	return w;
}

/*
 * Creates a new windowlist window
 */
pms_win_windowlist *	Display::create_windowlist()
{
	pms_win_windowlist *	w;

	w = new pms_win_windowlist(this, &playlists);
	if (w)
	{
		if (pms->options->topbar_lines.size() == 0 || !pms->options->topbarvisible)
			w->resize(0, 0, COLS, LINES - 1);
		else
			w->resize(0, pms->disp->topbar->height(), COLS, LINES - pms->disp->topbar->height() - 1);
		w->setborders(true, false, false, false);
		w->settitle("Windows");
		windows.push_back(w);
		if (curwin == NULL)
			curwin = w;
	}

	return w;
}

/*
 * Creates a new playlist window
 */
pms_win_playlist *	Display::create_playlist()
{
	pms_win_playlist *	w;
	
	w = new pms_win_playlist();
	if (w)
	{
		if (pms->options->topbar_lines.size() == 0 || !pms->options->topbarvisible)
			w->resize(0, 0, COLS, LINES - 1);
		else
			w->resize(0, pms->disp->topbar->height(), COLS, LINES - pms->disp->topbar->height() - 1);
		w->setborders(true, false, false, false);
		windows.push_back(w);
		playlists.push_back(w);
		if (curwin == NULL)
			curwin = w;
	}

	return w;
}

bool			Display::delete_window(pms_window * win)
{
	unsigned int			i;
	pms_window *			t;
	int				c = -1;
	bool				deleted = false;

	if (!win) return false;

	for (i = 0; i < windows.size(); i++)
	{
		t = windows[i];
		if (t == win)
		{
			if (t == curwin)
				c = static_cast<int>(i);
			delete windows[i];
			windows.erase(windows.begin() + i);
			deleted = true;
			break;
		}
	}
	for (i = 0; i < playlists.size(); i++)
	{
		if (dynamic_cast<pms_window *>(playlists[i]) == win)
		{
			playlists.erase(playlists.begin() + i);
			break;
		}
	}

	if (deleted && c >= 0)
	{
		i = static_cast<unsigned int>(c);
		curwin = NULL;
		if (i >= windows.size()) i = windows.size() - 1;
		if (curwin == NULL)
			activate(windows[i]);
	}

	return deleted;
}
