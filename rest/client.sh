#!/bin/bash

source "$(dirname "$0")/../.env"
get_all_clients() {
    curl -X GET "$BASE_URL/api/client/" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

get_clients_with_params() {
    local is_active=${1:-true}
    local page=${2:-1}
    curl -X GET "$BASE_URL/api/client/?is_active=$is_active&page=$page" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

get_clients_by_department() {
    local department_id=${1:--1}
    local page=${2:-1}
    curl -X GET "$BASE_URL/api/client/?departmentId=$department_id&page=$page" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

search_clients() {
    local department_id=${1:--1}
    local page=${2:-1}
    local per_page=${3:-400}
    curl -X GET "$BASE_URL/api/client/?departmentId=$department_id&page=$page&per_page=$per_page" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

search_clients_by_page_id() {
    local page=${1:-2}
    local page_id=${2:-2}
    curl -X GET "$BASE_URL/api/client/?page=$page&page_id=$page_id" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

case "${1:-all}" in
    "all")
        echo "Getting all clients..."
        get_all_clients
        ;;
    "params")
        echo "Getting clients with params (is_active=${2:-true}, page=${3:-1})..."
        get_clients_with_params "$2" "$3"
        ;;
    "department")
        echo "Getting clients by department (departmentId=${2:--1}, page=${3:-1})..."
        get_clients_by_department "$2" "$3"
        ;;
    "search")
        echo "Searching clients (departmentId=${2:--1}, page=${3:-1}, per_page=${4:-400})..."
        search_clients "$2" "$3" "$4"
        ;;
    "page")
        echo "Searching clients by page ID (page=${2:-2}, page_id=${3:-2})..."
        search_clients_by_page_id "$2" "$3"
        ;;
    *)
        echo "Usage: $0 [all|params|department|search|page] [args...]"
        echo "  all        - Get all clients (default)"
        echo "  params     - Get clients with params [is_active] [page]"
        echo "  department - Get clients by department [departmentId] [page]"
        echo "  search     - Search clients [departmentId] [page] [per_page]"
        echo "  page       - Search clients by page ID [page] [page_id]"
        ;;
esac