#!/bin/bash

set -o errexit
set -o nounset

readonly COMMAND={1:"none"}

main() {
	case $COMMAND in
	total)
		cat ./nodes.csv | wc -l
		;;
	total-per-country)
		cat ./nodes.csv | awk -F ',' '{print $$3}' | sort | uniq -c | sort
		;;
	reachable)
		cat ./nodes.csv | grep -v 'dial' | grep -v 'net' | grep -v 'reset' | wc -l
		;;
	reachable-per-country)
		cat ./nodes.csv | grep -v 'dial' | grep -v 'net' | grep -v 'reset' | awk -F ',' '{print $$3}' | sort | uniq -c | sort
		;;
	*)
		help $0
		exit 1
		;;
	esac
}

help () {
	echo "usage: $1 <command>
Commands:
	total
	total-per-country
	reachable
	reachable-per-country
	"
}

main "$@"
