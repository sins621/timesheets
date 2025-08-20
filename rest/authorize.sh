#!/bin/bash

source "$(dirname "$0")/../.env"
authorize() {
    local email=${1:-$EMAIL}
    local password=${2:-$PASSWORD}
    
    curl -X POST "$BASE_URL/api/account/Authorise" \
         -H "Content-Type: application/json" \
         -d "{\"Email\":\"$email\",\"Password\":\"$password\"}" \
         -s | jq '.'
}

case "${1:-auth}" in
    "auth")
        echo "Authorizing user..."
        authorize "$2" "$3"
        ;;
    *)
        echo "Usage: $0 auth [email] [password]"
        echo "  auth - Authorize user (uses EMAIL and PASSWORD from .env if not provided)"
        ;;
esac