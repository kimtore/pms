# adapted from https://github.com/SirCmpwn/sway

include(GNUInstallDirs)

add_custom_target(man ALL)

function(add_manpage filename name section)
	add_custom_command(
		OUTPUT ${name}.${section}
		COMMAND ${PANDOC}
				--standalone
				--to man
				${CMAKE_CURRENT_SOURCE_DIR}/${filename}
				-o ${name}.${section}
		DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/${filename}
		COMMENT Generating manpage for ${name}.${section}
	)

	add_custom_target(man-${name}.${section}
		DEPENDS
			${name}.${section}
	)
	add_dependencies(man
		man-${name}.${section}
	)

	install(
		FILES ${CMAKE_CURRENT_BINARY_DIR}/${name}.${section}
		DESTINATION ${CMAKE_INSTALL_MANDIR}/man${section}
		COMPONENT documentation
	)
endfunction()
