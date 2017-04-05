cmake_minimum_required(VERSION 3.3.1)

project(pms VERSION 0.99.0 LANGUAGES CXX)
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_EXTENSIONS OFF)

list(APPEND CMAKE_MODULE_PATH "${PROJECT_SOURCE_DIR}/cmake")

# Build the exeutable
include_directories("${PROJECT_SOURCE_DIR}/src")

add_executable(pms
   src/action.cpp
   src/color.cpp
   src/command.cpp
   src/config.cpp
   src/conn.cpp
   src/display.cpp
   src/error.cpp
   src/field.cpp
   src/input.cpp
   src/library.cpp
   src/list.cpp
   src/listitem.cpp
   src/message.cpp
   src/options.cpp
   src/playlist.cpp
   src/pms.cpp
   src/queue.cpp
   src/search.cpp
   src/set_parameters.cpp
   src/song.cpp
   src/songlist.cpp
)

set(CURSES_NEED_WIDE true)

include(FindPkgConfig)
find_package(Threads REQUIRED)
find_package(Ncursesw REQUIRED)
pkg_check_modules(GLIB REQUIRED glib-2.0>=2.0)
pkg_check_modules(MPDC REQUIRED libmpdclient>=2.5)

option(regex "Enable regex support" YES)
if(regex)
	if(CMAKE_COMPILER_IS_GNUCXX)
		execute_process(COMMAND ${CMAKE_CXX_COMPILER} -dumpversion OUTPUT_VARIABLE GXX_VERSION)
		if(GXX_VERSION VERSION_LESS 4.9)
			message("g++ of at least 4.9 is required for proper regex support.")
			set(regex NO)
		endif()
	endif()
endif()

if(regex)
	target_compile_definitions(pms PUBLIC -DHAVE_REGEX=1)
endif()

target_compile_definitions(pms PUBLIC -DCMAKE_BUILD)

target_include_directories(pms PUBLIC
	${CURSES_INCLUDE_PATH}
	${GLIB_INCLUDE_DIRS}
	${MPDC_INCLUDE_DIRS}
)

target_link_libraries(pms
	${CMAKE_THREADS_LIBS_INIT}
	${CURSES_LIBRARIES}
	${GLIB_LIBRARIES}
	${MPDC_LIBRARIES}
)

include(GNUInstallDirs)
install(TARGETS pms
	RUNTIME
	DESTINATION ${CMAKE_INSTALL_BINDIR}
	COMPONENT runtime
)

# Build the manpage
option(ENABLE_DOC "Build documentation" YES)
find_program(PANDOC pandoc)
if(NOT PANDOC)
	set(ENABLE_DOC NO)
endif()
if(ENABLE_DOC)
	include(Manpage)
	add_manpage(doc/pms.md pms 1)
endif()
