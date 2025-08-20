#!/bin/bash

source "$(dirname "$0")/../.env"
get_all_projects() {
    curl -X GET "$BASE_URL/api/Project/" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

get_project_by_id() {
    local project_id=${1:-1}
    curl -X GET "$BASE_URL/api/Project/$project_id" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

case "${1:-all}" in
    "all")
        echo "Getting all projects..."
        get_all_projects
        ;;
    "id")
        echo "Getting project by ID: ${2:-1}..."
        get_project_by_id "${2:-1}"
        ;;
    *)
        echo "Usage: $0 [all|id] [project_id]"
        echo "  all - Get all projects (default)"
        echo "  id  - Get project by ID (default: 1)"
        ;;
esac