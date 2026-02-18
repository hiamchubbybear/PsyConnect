#!/bin/bash

# Define keywords to alert on
KEYWORDS="password=|key=|secret=|token="

# Find all tracked .env* files (env, env.example, env.dev, etc)
# Exclude node_modules and vendor
TRACKED_CONFIG_FILES=$(git ls-files | grep -E "\.env(\..+)?$" | grep -v "node_modules" | grep -v "vendor")

if [ -z "$TRACKED_CONFIG_FILES" ]; then
    echo "✅ No config files found to check."
    exit 0
fi

echo "🔍 Scanning tracked config files for exposed secrets..."
has_error=0

for file in $TRACKED_CONFIG_FILES; do
    # Check if file contains suspicious assignments (KEY=value where value is not empty or placeholder)
    # We grep for lines that have keywords followed by something that isn't a placeholder like 'your_', 'Insert', '...'
    if grep -E "($KEYWORDS)" "$file" | grep -vE "your_|YOUR_|placeholder|INSERT|EXAMPLE" > /dev/null; then
        echo "❌ POTENTIAL SECRET LEAK IN: $file"
        grep -E "($KEYWORDS)" "$file" | grep -vE "your_|YOUR_|placeholder|INSERT|EXAMPLE"
        has_error=1
    fi
done

if [ $has_error -eq 1 ]; then
    echo ""
    echo "⚠️  SECURITY CHECK FAILED!"
    echo "Please replace actual credentials with placeholders (e.g., your_password_here) before committing."
    echo "Use 'git commit --no-verify' ONLY if you are absolutely sure this is a false positive."
    exit 1
fi

echo "✅ Security Check Passed."
exit 0
