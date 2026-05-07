#!/bin/bash
# Fetch Stitch design assets with proper redirect handling
# Usage: ./fetch-stitch.sh <url> <output_path>

set -e

URL="$1"
OUTPUT="$2"

if [ -z "$URL" ] || [ -z "$OUTPUT" ]; then
    echo "Usage: $0 <url> <output_path>"
    exit 1
fi

echo "Fetching: $URL"
echo "Output: $OUTPUT"

# Create directory if needed
mkdir -p "$(dirname "$OUTPUT")"

# Fetch with redirect handling and proper headers
curl -L -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
     -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" \
     -o "$OUTPUT" \
     "$URL"

echo "Downloaded to: $OUTPUT"
