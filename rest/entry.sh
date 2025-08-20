#!/bin/bash

source "$(dirname "$0")/../.env"
get_all_entries() {
    curl -X GET "$BASE_URL/api/entry/" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

get_entry_by_id() {
    local entry_id=${1:-1}
    curl -X GET "$BASE_URL/api/entry/$entry_id" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq '.'
}

create_entry() {
    local comments=${1:-"test"}
    local entry_date=${2:-$(date +%Y-%m-%d)}
    local time=${3:-8}
    local overtime=${4:-0}
    local person_id=${5:-123}
    local task_id=${6:-123}
    local cost_code_id=${7:-1}
    
    curl -X POST "$BASE_URL/api/entry/create" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"Comments\": \"$comments\",
            \"EntryDate\": \"$entry_date\",
            \"Time\": $time,
            \"Overtime\": $overtime,
            \"Person\": {\"PersonId\": $person_id},
            \"Task\": {\"TaskId\": $task_id},
            \"CostCodeId\": $cost_code_id
        }" \
        -s | jq '.'
}

find_user_by_name() {
    local name=${1:-"Bradly Carpenter"}
    echo "Searching for user: $name..."
    curl -X GET "$BASE_URL/api/users/" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq --arg name "$name" '.[] | select(.name // .fullName // .firstName + " " + .lastName | test($name; "i"))'
}

find_project_by_name() {
    local project_name=${1:-"Lumix"}
    echo "Searching for project: $project_name..."
    curl -X GET "$BASE_URL/api/Project/" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -H "Content-Type: application/json" \
        -s | jq --arg name "$project_name" '.[] | select(.name // .projectName | test($name; "i"))'
}

upload_attachment() {
    local entry_id=${1:-1}
    local file_path=${2:-"./example-file.txt"}
    
    if [[ ! -f "$file_path" ]]; then
        echo "Error: File '$file_path' not found!"
        return 1
    fi
    
    curl -X POST "$BASE_URL/api/entry/$entry_id/uploadattachment" \
        -H "Authorization: Bearer $BEARER_TOKEN" \
        -F "file=@$file_path" \
        -s | jq '.'
}

case "${1:-all}" in
    "all")
        echo "Getting all entries..."
        get_all_entries
        ;;
    "id")
        echo "Getting entry by ID: ${2:-1}..."
        get_entry_by_id "${2:-1}"
        ;;
    "create")
        echo "Creating new entry..."
        echo "Usage: create [comments] [entry_date] [time] [overtime] [person_id] [task_id] [cost_code_id]"
        create_entry "$2" "$3" "$4" "$5" "$6" "$7" "$8"
        ;;
    "find-user")
        echo "Finding user by name: ${2:-Bradly Carpenter}..."
        find_user_by_name "$2"
        ;;
    "find-project")
        echo "Finding project by name: ${2:-Lumix}..."
        find_project_by_name "$2"
        ;;
    "create-for-bradly-lumix")
        echo "Creating entry for Bradly Carpenter on Lumix project..."
        echo "First, let's find your PersonId..."
        USER_INFO=$(find_user_by_name "Bradly Carpenter")
        echo "User info: $USER_INFO"
        
        echo "Now finding Lumix project..."
        PROJECT_INFO=$(find_project_by_name "Lumix")
        echo "Project info: $PROJECT_INFO"
        
        create_entry "test" "$(date +%Y-%m-%d)" "8" "0" "123" "123" "1"
        ;;
    "upload")
        echo "Uploading attachment to entry ${2:-1} from file ${3:-./example-file.txt}..."
        upload_attachment "$2" "$3"
        ;;
    *)
        echo "Usage: $0 [all|id|create|find-user|find-project|create-for-bradly-lumix|upload] [args...]"
        echo "  all                    - Get all entries (default)"
        echo "  id                     - Get entry by ID [entry_id]"
        echo "  create                 - Create new entry [comments] [entry_date] [time] [overtime] [person_id] [task_id] [cost_code_id]"
        echo "  find-user             - Find user by name [name]"
        echo "  find-project          - Find project by name [project_name]"
        echo "  create-for-bradly-lumix - Create entry for Bradly Carpenter on Lumix project"
        echo "  upload                - Upload attachment [entry_id] [file_path]"
        ;;
esac