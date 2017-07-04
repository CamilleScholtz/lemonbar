#!/bin/sh

go build
./bar | ../lemonbar \
	-g 1908x18 \
	-r 6 \
	-B '#423F31' -R '#423F31' \
	-f '-jmk-neep-medium-r-normal--10-80-75-75-c-50-iso8859-1' \
	-f '-wuncon-siji-medium-r-normal--10-100-75-75-c-80-iso10646-1' \
	-f '-freetype-kochismall-medium-r-normal--9-90-75-75-p-77-iso10646-1' \
	-f '-freetype-baekmuksmall-medium-r-normal--9-90-75-75-p-38-iso10646-1'
