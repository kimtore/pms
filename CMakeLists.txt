cmake_minimum_required(VERSION 2.8)
project("Practical Music Search" CXX)
set(PMS_VERSION_MAJOR 0)
set(PMS_VERSION_MINOR 50)
add_executable(pms src/main.cpp src/curses.cpp src/config.cpp src/console.cpp src/mpd.cpp src/wconsole.cpp src/window.cpp src/wstatusbar.cpp src/windowmanager.cpp src/input.cpp src/command.cpp src/keybinding.cpp src/pms.cpp src/color.cpp src/songlist.cpp src/song.cpp src/wsonglist.cpp src/field.cpp src/wtopbar.cpp src/topbar.cpp src/search.cpp src/sort.cpp src/clipboard.cpp)
set(CMAKE_CXX_FLAGS, "-Wall -Werror -g3 -pedantic")

configure_file(
	"${PROJECT_SOURCE_DIR}/src/build.h.in"
	"${PROJECT_BINARY_DIR}/src/build.h"
)

set(CURSES_NEED_NCURSES, true)
include(FindCurses)
if (NOT CURSES_FOUND)
	message(FATAL_ERROR "PMS requires NCurses >= 5.0.")
endif(NOT CURSES_FOUND)
include_directories(${CURSES_INCLUDE_DIR})
target_link_libraries(pms ${CURSES_LIBRARIES})
