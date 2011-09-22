cmake_minimum_required(VERSION 2.8)
project("Practical Music Search" CXX)
set(PMS_VERSION_MAJOR 0)
set(PMS_VERSION_MINOR 50)
add_executable(pms src/main.cpp src/curses.cpp src/config.cpp src/console.cpp src/mpd.cpp src/wconsole.cpp src/window.cpp src/wstatusbar.cpp src/windowmanager.cpp src/input.cpp src/command.cpp src/keybinding.cpp src/pms.cpp)
set(CMAKE_CXX_FLAGS, "-Wall -g -pedantic")

configure_file(
	"${PROJECT_SOURCE_DIR}/src/build.h.in"
	"${PROJECT_BINARY_DIR}/src/build.h"
)

include(FindGLIB2.cmake)
if (NOT GLIB2_FOUND)
	message(FATAL_ERROR "PMS requires GLIB 2.0 to compile.")
endif(NOT GLIB2_FOUND)	
include_directories(${GLIB2_INCLUDE_DIRS})
target_link_libraries(pms ${GLIB2_LIBRARIES})

set(CURSES_NEED_NCURSES, true)
include(FindCurses.cmake)
if (NOT CURSES_FOUND)
	message(FATAL_ERROR "PMS requires ncurses to compile.")
endif(NOT CURSES_FOUND)
include_directories(${CURSES_INCLUDE_DIR})
target_link_libraries(pms ${CURSES_LIBRARIES})
