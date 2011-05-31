cmake_minimum_required(VERSION 2.8)
project("Practical Music Search" C CXX)
set(PMS_VERSION_MAJOR 0)
set(PMS_VERSION_MINOR 50)
add_executable(pms src/main.cpp)

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
