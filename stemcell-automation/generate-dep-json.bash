#!/usr/bin/env bash

set -euo pipefail

sha_open_ssh=$(shasum -a 256 $1 | cut -d " " -f 1)
sha_bosh_psmodules=$(shasum -a 256 $2 | cut -d " " -f 1)
sha_agent=$(shasum -a 256 $3 | cut -d " " -f 1)
sha_lgpo=$(shasum -a 256 $4 | cut -d " " -f 1)

cat <<EOF
{
  "OpenSSH-Win64.zip": {
    "sha": "${sha_open_ssh}"
  },
  "bosh-psmodules.zip": {
    "sha": "${sha_bosh_psmodules}"
  },
  "agent.zip": {
    "sha": "${sha_agent}"
  },
  "LGPO.zip": {
    "sha": "${sha_lgpo}"
  }
}
EOF
