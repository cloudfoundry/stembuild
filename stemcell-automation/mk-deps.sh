#!/bin/bash
set -e

DEPS=(
  agent.zip
  LGPO.zip
  OpenSSH-Win64.zip
  bosh-psmodules.zip
)

# check that we have all these files on disk?

ruby mk-deps.rb "${DEPS[@]}" | jq . > deps2.json
