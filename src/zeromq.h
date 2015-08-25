/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS	<<Practical Music Search>>
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


#ifndef _ZEROMQ_H_
#define _ZEROMQ_H_

#include <zmq.h>
#include <pthread.h>
#include <mpd/client.h>

#define ZEROMQ_SOCKET_IDLE "inproc://idle"
#define ZEROMQ_SOCKET_INPUT "inproc://input"

/**
 * ZeroMQ inter-thread communication
 */
class ZeroMQ
{
private:
	void *				context;
	void *				socket_idle;
	void *				socket_input;
	zmq_pollitem_t			poll_items[2];

	pthread_t			idle_thread;
	pthread_t			input_thread;

public:
					ZeroMQ();

	bool				has_idle_events();
	enum mpd_idle			get_idle_events();
	void				continue_idle();
	bool				has_input_events();
	wchar_t				get_input_events();
	void				start_thread_idle(void *(*func) (void *));
	void				start_thread_input(void *(*func) (void *));
	void				poll_events(int timeout);
};

#endif /* _ZEROMQ_H_ */
