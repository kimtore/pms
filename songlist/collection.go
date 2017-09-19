package songlist

import (
	"fmt"
	"reflect"
	"time"
)

// Collection holds a set of songlists, and keeps track of movement between different lists.
type Collection struct {
	lists   []Songlist
	index   int
	last    Songlist
	current Songlist
	updated time.Time
}

// NewCollection returns Collection.
func NewCollection() *Collection {
	return &Collection{
		lists: make([]Songlist, 0),
	}
}

// Activate activates the specified songlist. If the songlist is indexed, the
// index is switched to that of the songlist.
func (c *Collection) Activate(s Songlist) {
	c.index = -1
	for i, stored := range c.lists {
		if stored == s {
			c.index = i
			break
		}
	}
	c.last = c.current
	c.current = s
	c.setUpdated()
}

func (c *Collection) ActivateIndex(i int) error {
	list, err := c.Songlist(i)
	if err != nil {
		return err
	}
	c.Activate(list)
	return nil
}

// Add appends a songlist to the collection.
func (c *Collection) Add(s Songlist) {
	c.lists = append(c.lists, s)
}

func (c *Collection) Current() Songlist {
	return c.current
}

func (c *Collection) Index() (int, error) {
	if !c.ValidIndex(c.index) {
		return 0, fmt.Errorf("Songlist index is out of range")
	}
	return c.index, nil
}

// Last returns the last used songlist.
func (c *Collection) Last() Songlist {
	return c.last
}

// Len returns the songlists count.
func (c *Collection) Len() int {
	return len(c.lists)
}

func (c *Collection) Remove(index int) error {
	if err := c.ValidateIndex(index); err != nil {
		return err
	}
	if index+1 == c.Len() {
		c.lists = c.lists[:index]
	} else {
		c.lists = append(c.lists[:index], c.lists[index+1:]...)
	}
	return nil
}

// Replace replaces an existing songlist with its new version. Checking
// is done on a type-level, so only the queue and library will be replaced.
func (c *Collection) Replace(s Songlist) {
	for i := range c.lists {
		if reflect.TypeOf(c.lists[i]) != reflect.TypeOf(s) {
			continue
		}
		//console.Log("Songlist UI: replacing songlist of type %T at %p with new list at %p", s, c.lists[i], s)
		//console.Log("Songlist UI: comparing %p %p", c.lists[i], c.Songlist())

		active := c.lists[i] == c.Current()
		c.lists[i] = s

		if active {
			//console.Log("Songlist UI: replaced songlist is currently active, switching to new songlist.")
			c.Activate(s)
		}
		return
	}

	//console.Log("Songlist UI: adding songlist of type %T at address %p since no similar exists", s, s)
	c.Add(s)
}

// newVersion updates the version number of this data set.
func (c *Collection) setUpdated() {
	c.updated = time.Now()
}

func (c *Collection) Songlist(index int) (Songlist, error) {
	if err := c.ValidateIndex(index); err != nil {
		return nil, err
	}
	return c.lists[index], nil
}

func (c *Collection) ValidIndex(i int) bool {
	return i >= 0 && i < c.Len()
}

func (c *Collection) ValidateIndex(i int) error {
	if !c.ValidIndex(i) {
		return fmt.Errorf("Index %d is out of bounds (try between 1 and %d)", i+1, c.Len())
	}
	return nil
}

func (c *Collection) Updated() time.Time {
	return c.updated
}
