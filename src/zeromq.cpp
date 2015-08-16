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

#include "zeromq.h"

#include <zmq.h>
#include <assert.h>
#include <stdlib.h>
#include <pthread.h>
#include <mpd/client.h>

/**
 * Poll the input and IDLE subsystems for events, and block for up to `timeout'
 * milliseconds, or or until either MPD or the user makes some noise.
 */
void
ZeroMQ::poll_events(int timeout)
{
	while(true) {
		if (zmq_poll(poll_items, 2, timeout) == -1) {
			if (errno == EINTR) {
				continue;
			}
			abort();
		}
		return;
	}
}

/* Check for events on the IDLE socket. */
bool
ZeroMQ::has_idle_events()
{
	return (poll_items[0].revents & ZMQ_POLLIN);
}

/* Retrieve events from the IDLE socket. */
enum mpd_idle
ZeroMQ::get_idle_events()
{
	int rc;
	enum mpd_idle idle_reply;

	rc = zmq_recv(socket_idle, (void *)&idle_reply, sizeof(enum mpd_idle *), 0);
	assert(rc == sizeof(enum mpd_idle *));

	return idle_reply;
}

/* Send a signal to the IDLE thread to continue polling. */
void
ZeroMQ::continue_idle()
{
	assert(zmq_send(socket_idle, NULL, 0, 0) == 0);
}

/* Check for events on the input socket. */
bool
ZeroMQ::has_input_events()
{
	return (poll_items[1].revents & ZMQ_POLLIN);
}

/* Retrieve events from the input socket. */
wchar_t
ZeroMQ::get_input_events()
{
	int rc;
	wchar_t reply;

	rc = zmq_recv(socket_input, (void *)&reply, sizeof(wchar_t *), 0);
	assert(rc == sizeof(wchar_t *));

	return reply;
}

/**
 * Initialize the MPD IDLE thread.
 */
void
ZeroMQ::start_thread_idle(void *(*func) (void *))
{
	assert(pthread_create(&idle_thread, NULL, func, context) == 0);
}

/**
 * Initialize the ncurses input thread.
 */
void
ZeroMQ::start_thread_input(void *(*func) (void *))
{
	assert(pthread_create(&input_thread, NULL, func, context) == 0);
}

/**
 * Initialize ZeroMQ and the main thread's ZeroMQ REQ/SUB sockets.
 */
ZeroMQ::ZeroMQ()
{
	/* Initialize ZeroMQ context and sockets */
	context = zmq_ctx_new();
	assert(context != NULL);
	socket_idle = zmq_socket(context, ZMQ_REQ);
	assert(socket_idle != NULL);
	socket_input = zmq_socket(context, ZMQ_SUB);
	assert(socket_input != NULL);
	assert(zmq_setsockopt(socket_input, ZMQ_SUBSCRIBE, "", 0) == 0);
	assert(zmq_connect(socket_idle, ZEROMQ_SOCKET_IDLE) == 0);
	assert(zmq_bind(socket_input, ZEROMQ_SOCKET_INPUT) == 0);

	/* Set up ZeroMQ poller */
	poll_items[0].socket = socket_idle;
	poll_items[0].events = ZMQ_POLLIN;
	poll_items[1].socket = socket_input;
	poll_items[1].events = ZMQ_POLLIN;
}
