/*
 * T er for eksempel Song, Window, Object, alt mulig
 */

template <class T>
class ListView
{
	T *		list;
public:
	void		select(long pos, int state);
}


class Window
{
public:
	int			cursor_x;
	int			cursor_y;
}
