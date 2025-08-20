#!/bin/bash

source "$(dirname "$0")/../.env"
login() {
    local email=${1:-$EMAIL}
    local password=${2:-$PASSWORD}
    
    echo "Logging in as $email..."
    
    response=$(curl -X POST "$BASE_URL/api/account/Authorise" \
                    -H "Content-Type: application/json" \
                    -d "{\"Email\":\"$email\",\"Password\":\"$password\"}" \
                    -s)
    
    if [[ $? -ne 0 ]]; then
        echo "Error: Failed to connect to server"
        return 1
    fi
    
    token=$(echo "$response" | jq -r '.token // .access_token // .bearerToken // .authToken // empty')
    
    if [[ -z "$token" || "$token" == "null" ]]; then
        echo "Error: No token found in response"
        echo "Response: $response"
        return 1
    fi
    
    echo "Login successful! Token: ${token:0:20}..."
    
    env_file="$(dirname "$0")/../.env"
    
    if grep -q "^BEARER_TOKEN=" "$env_file"; then
        sed -i "s/^BEARER_TOKEN=.*/BEARER_TOKEN=$token/" "$env_file"
    else
        echo "BEARER_TOKEN=$token" >> "$env_file"
    fi
    
    echo "Updated BEARER_TOKEN in .env file"
    
    export BEARER_TOKEN="$token"
    
    return 0
}

case "${1:-login}" in
    "login")
        login "$2" "$3"
        ;;
    *)
        echo "Usage: $0 login [email] [password]"
        echo "  login - Login and update token (uses EMAIL and PASSWORD from .env if not provided)"
        ;;
esac