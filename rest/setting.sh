#!/bin/bash

source "$(dirname "$0")/../.env"
get_setting_by_key() {
    local key=${1:-example-key}
    curl -X GET "$BASE_URL/api/setting/$key" \
         -H "Authorization: Bearer $BEARER_TOKEN" \
         -H "Content-Type: application/json" \
         -s | jq '.'
}

case "${1:-key}" in
    "key")
        echo "Getting setting by key: ${2:-example-key}..."
        get_setting_by_key "${2:-example-key}"
        ;;
    *)
        echo "Usage: $0 key [setting_key]"
        echo "  key - Get setting by key (default: example-key)"
        ;;
esac