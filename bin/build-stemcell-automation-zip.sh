#!/usr/bin/env bash
set -eu -o pipefail

ROOT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
STEMCELL_AUTOMATION_PS1=$(ls "${ROOT_DIR}"/stemcell-automation/*ps1 | grep -iv Test)

: ${OPENSSH_ZIP?"Please see README.md on where to obtain this."}
: ${BOSH_PSMODULES_ZIP?"Please see README.md on where to obtain this."}
: ${AGENT_ZIP?"Please see README.md on how to construct this."}
: ${DEPS_JSON?"Please see README.md on how to construct this."}

TEMP_DIR=$(mktemp -d)

cp "${OPENSSH_ZIP}" "${TEMP_DIR}/OpenSSH-Win64.zip"
cp "${BOSH_PSMODULES_ZIP}" "${TEMP_DIR}/bosh-psmodules.zip"
cp "${AGENT_ZIP}" "${TEMP_DIR}/agent.zip"
cp "${DEPS_JSON}" "${TEMP_DIR}/deps.json"
cp ${STEMCELL_AUTOMATION_PS1} "$TEMP_DIR"

rm -f "${ROOT_DIR}/assets/StemcellAutomation.zip"

zip -rj "${ROOT_DIR}/assets/StemcellAutomation.zip" "$TEMP_DIR"

rm -r "$TEMP_DIR"