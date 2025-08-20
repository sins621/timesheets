#!/bin/bash

source "$(dirname "$0")/../.env"
get_all_users() {
    curl -X GET "$BASE_URL/api/users/?page=1&per_page=100" \
         -H "Authorization: Bearer $BEARER_TOKEN" \
         -H "Content-Type: application/json" \
         -s | jq '.'
}

get_current_user() {
    curl -X GET "$BASE_URL/api/users/me" \
         -H "Authorization: Bearer $BEARER_TOKEN" \
         -H "Content-Type: application/json" \
         -s | jq '.'
}

get_user_by_id() {
    local user_id=${1:-1}
    curl -X GET "$BASE_URL/api/users/$user_id" \
         -H "Authorization: Bearer $BEARER_TOKEN" \
         -H "Content-Type: application/json" \
         -s | jq '.'
}

case "${1:-all}" in
    "all")
        echo "Getting all users..."
        get_all_users
        ;;
    "me")
        echo "Getting current user..."
        get_current_user
        ;;
    "id")
        echo "Getting user by ID: ${2:-1}..."
        get_user_by_id "${2:-1}"
        ;;
    *)
        echo "Usage: $0 [all|me|id] [user_id]"
        echo "  all - Get all users (default)"
        echo "  me  - Get current user"
        echo "  id  - Get user by ID (default: 1)"
        ;;
esac