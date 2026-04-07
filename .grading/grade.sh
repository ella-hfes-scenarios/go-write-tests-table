#!/bin/bash
set -euo pipefail
NONCE=$(openssl rand -hex 16)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORK_DIR=$(mktemp -d)
trap "rm -rf $WORK_DIR" EXIT

# Copy candidate's test files and go.mod
find . -name '*_test.go' -exec cp --parents {} "$WORK_DIR/" \;
cp go.mod go.sum "$WORK_DIR/" 2>/dev/null || true

# Run against reference
cp -r "$SCRIPT_DIR/reference/pkg/" "$WORK_DIR/pkg/"
cd "$WORK_DIR"
REF_OUTPUT=$(go test ./... -count=1 -json 2>&1 || true)
REF_PASSED=$(echo "$REF_OUTPUT" | grep '"Test"' | grep '"pass"' | wc -l | tr -d ' ')
REF_TOTAL=$(echo "$REF_OUTPUT" | grep '"Test"' | grep -E '"pass"|"fail"' | wc -l | tr -d ' ')
REF_RATE=$(echo "scale=2; ${REF_PASSED:-0} / (${REF_TOTAL:-1})" | bc)

MUTATIONS_CAUGHT=0
MUTATIONS_TOTAL=0
MUTATION_RESULTS="["
for buggy_dir in "$SCRIPT_DIR"/buggy-v*/; do
  variant=$(basename "$buggy_dir")
  rm -rf "$WORK_DIR/pkg/"
  cp -r "$buggy_dir/pkg/" "$WORK_DIR/pkg/"
  OUTPUT=$(cd "$WORK_DIR" && go test ./... -count=1 2>&1 || true)
  if echo "$OUTPUT" | grep -q "FAIL"; then
    CAUGHT="true"
    MUTATIONS_CAUGHT=$((MUTATIONS_CAUGHT + 1))
  else
    CAUGHT="false"
  fi
  MUTATIONS_TOTAL=$((MUTATIONS_TOTAL + 1))
  MUTATION_RESULTS="$MUTATION_RESULTS{\"variant\":\"$variant\",\"caught\":$CAUGHT},"
done
MUTATION_RESULTS="${MUTATION_RESULTS%,}]"
MUTATION_SCORE=$(echo "scale=2; $MUTATIONS_CAUGHT / ($MUTATIONS_TOTAL + 0.001)" | bc)

echo "{\"referencePassRate\":$REF_RATE,\"mutationResults\":$MUTATION_RESULTS,\"mutationsCaught\":$MUTATIONS_CAUGHT,\"mutationsTotal\":$MUTATIONS_TOTAL,\"mutationScore\":$MUTATION_SCORE,\"nonce\":\"$NONCE\"}"
