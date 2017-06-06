# Configuring PMS and MPD

When starting the program, PMS connects to the MPD server specified in the `$MPD_HOST` and `$MPD_PORT` environment variables.

In order to create a full-text search index for fast searches, PMS retrieves the entire song library from MPD whenever the library is updated, and on every startup. If your song library is big, the `listallinfo` command will overflow MPD's send buffer, and the connection is dropped. This can be mitigated by increasing MPD's output buffer size, and then restarting MPD:

```
cat >>/etc/mpd.conf<<<EOF
max_output_buffer_size "262144"
EOF
/etc/init.d/mpd restart  # or equivalent
```
