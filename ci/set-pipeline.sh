#!/bin/bash

set -eo pipefail

TARGET=pcf-pipelines
PIPELINE=gershman-ptracker-resource

case "$1" in
destroy)
	fly -t "$TARGET" \
		destroy-pipeline -p "$PIPELINE"
	;;
*)
	fly -t "$TARGET" \
		set-pipeline -p "$PIPELINE" \
		-c pipeline.yml \
		-l ../scripts/credentials.yml \
		"$@"
	;;
esac
