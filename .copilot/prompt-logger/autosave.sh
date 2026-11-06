#!/bin/bash
# Auto-save prompts to project folder
# Usage: ./autosave.sh "prompt content"

PROMPT_DIR=".copilot/prompt-logger"
TIMESTAMP=$(date +"%Y-%m-%d-%H%M%S")
PROMPT_FILE="${PROMPT_DIR}/prompt-${TIMESTAMP}.md"

# Create directory if not exists
mkdir -p "${PROMPT_DIR}"

# Create markdown file with metadata
cat > "${PROMPT_FILE}" << EOF
---
timestamp: $(date -Iseconds)
project: $(basename "$PWD")
directory: $PWD
commit: $(git rev-parse --short HEAD 2>/dev/null || echo "not-a-git-repo")
---

# Prompt: ${TIMESTAMP}

**Timestamp**: $(date "+%Y-%m-%d %H:%M:%S")
**Directory**: \`$PWD\`
**Project**: \`$(basename "$PWD")\`

## Content

${1:-No content provided}

## Notes

Auto-saved by Copilot Prompt Logger
EOF

echo "Prompt saved to: ${PROMPT_FILE}"
