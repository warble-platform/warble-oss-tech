#!/bin/bash

# Make sure flarectl is in the PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Check if CF_API_TOKEN is set
if [ -z "$CF_API_TOKEN" ]; then
    echo "Error: CF_API_TOKEN environment variable is not set."
    echo "Please set it using: export CF_API_TOKEN='your_api_token'"
    exit 1
fi

DOMAINS=(
    "avirka.ai"
    "frakma.io"
    "warblecloud.ai"
    "warblecloud.com"
    "warblelabs.io"
    "warbleoss.org"
)

echo "Establishing security controls for Cloudflare domains..."

for domain in "${DOMAINS[@]}"; do
    echo "----------------------------------------"
    echo "Configuring security for: $domain"
    echo "----------------------------------------"

    # Set Security Level to High
    echo " -> Setting Security Level to High..."
    flarectl zone settings update --zone "$domain" --setting "security_level" --value "high"
    
    # Enable Always Use HTTPS
    echo " -> Enabling Always Use HTTPS..."
    flarectl zone settings update --zone "$domain" --setting "always_use_https" --value "on"
    
    # Set Minimum TLS Version to 1.2
    echo " -> Setting Minimum TLS Version to 1.2..."
    flarectl zone settings update --zone "$domain" --setting "min_tls_version" --value "1.2"
    
    # Enable Browser Integrity Check
    echo " -> Enabling Browser Integrity Check..."
    flarectl zone settings update --zone "$domain" --setting "browser_check" --value "on"
    
    # Enable Automatic HTTPS Rewrites
    echo " -> Enabling Automatic HTTPS Rewrites..."
    flarectl zone settings update --zone "$domain" --setting "automatic_https_rewrites" --value "on"
    
    echo "✅ Finished configuring $domain"
done

echo "========================================"
echo "All domains have been configured!"
echo "========================================"
