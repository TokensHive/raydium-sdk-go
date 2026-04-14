#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

tmp_profile="$(mktemp)"
trap 'rm -f "${tmp_profile}"' EXIT

go test ./... -coverprofile="${tmp_profile}" >/dev/null

failed=0
while IFS= read -r line; do
  if [[ "${line}" =~ ^total: ]]; then
    continue
  fi

  coverage="$(awk '{print $3}' <<<"${line}")"
  package_path="$(awk '{print $1}' <<<"${line}")"

  if [[ "${coverage}" == "[no" ]]; then
    continue
  fi

  if [[ "${coverage}" != "100.0%" ]]; then
    echo "coverage gate failed: ${package_path} => ${coverage}" >&2
    failed=1
  fi
done < <(go tool cover -func="${tmp_profile}")

if [[ "${failed}" -ne 0 ]]; then
  exit 1
fi

echo "coverage gate passed: all packages at 100.0%"
