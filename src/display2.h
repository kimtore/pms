#include <stdio.h>
#include <vector>

using namespace std;

typedef list_t			long;


/*
 * Property
 */
class DispListProp
{
public:
	bool			selected;
};


/*
 * Basic list, virtualizes all functionality
 */
template <class T>
class List
{
protected:
	vector<T>		items;

public:
	virtual list_t		match() {};
};


/*
 * Display list
 */
template <class T>
class DispList : protected List <T>
{
protected:

	vector<DispListProp *>		props;
public:

	T				item(int i) { return items[i]; };

	Prop *				prop(int i) { return props[i]; };
	int				add(T item)
	{
		items.push_back(item);
		props.push_back(new Prop);

		return items.size() - 1;
	}

	int				cursorpos;
};

class Song
{
public:
	int			title;
};




int main()
{
	List<Song *>		songs;
	List<int>		ints;
	Song *			mysong;
	int			pos;

	mysong = new Song();
	mysong->title = 10;

	pos = songs.add(mysong);
	songs.prop(pos)->selected = true;

	printf("Title has the value: %d\n", songs.item(pos)->title);


	pos = ints.add(3);
	songs.prop(pos)->selected = true;

	printf("Int has the value: %d\n", ints.item(pos)->title);
}
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
