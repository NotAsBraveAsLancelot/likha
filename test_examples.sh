#!/bin/bash

# Test script for Likha examples
# Runs likha with each example config and displays the output

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_header() {
    echo -e "\n${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Check if likha binary exists
if [ ! -f "./likha" ]; then
    print_error "likha binary not found in current directory"
    echo "Please build the project first: go build -o likha"
    exit 1
fi

# Create examples directory if it doesn't exist
mkdir -p examples

# Array of example files (excluding custom ones)
example_files=(
    "basic_users.yaml"
    "ecommerce_products.yaml"
    "server_logs.yaml"
    "financial_transactions.yaml"
    "iot_sensors.yaml"
    "social_posts.yaml"
    "complex_csv.yaml"
    "xml_records.yaml"
    "yaml_config.yaml"
)

# Main execution
print_header "Testing Likha Examples"
echo "Generating 5 records for each example configuration..."

for config_file in "${example_files[@]}"; do
    config_path="examples/$config_file"
    
    # Skip if file doesn't exist
    if [ ! -f "$config_path" ]; then
        print_warning "Config file $config_path not found, skipping..."
        continue
    fi
    
    print_header "Testing: $config_file"
    
    # Run likha with the config
    echo "Running: ./likha --config $config_path --count 5"
    
    if ./likha --config "$config_path" --count 5; then
        print_success "Successfully generated data for $config_file"
        
        # Determine output file name from config
        case $config_file in
            "basic_users.yaml")
                output_file="basic_users.json"
                ;;
            "ecommerce_products.yaml")
                output_file="ecommerce_products.json"
                ;;
            "server_logs.yaml")
                output_file="server_logs.json"
                ;;
            "financial_transactions.yaml")
                output_file="financial_transactions.json"
                ;;
            "iot_sensors.yaml")
                output_file="iot_sensors.json"
                ;;
            "social_posts.yaml")
                output_file="social_posts.json"
                ;;
            "complex_csv.yaml")
                output_file="complex_data.csv"
                ;;
            "xml_records.yaml")
                output_file="xml_records.xml"
                ;;
            "yaml_config.yaml")
                output_file="yaml_config.yaml"
                ;;
            *)
                print_warning "Unknown output file for $config_file"
                continue
                ;;
        esac
        
        # Display the generated output
        if [ -f "$output_file" ]; then
            echo -e "\n${YELLOW}Generated output:${NC}"
            echo "----------------------------------------"
            cat "$output_file"
            echo "----------------------------------------"
            
            # Clean up output file
            rm -f "$output_file"
        else
            print_warning "Output file $output_file not found"
        fi
        
    else
        print_error "Failed to generate data for $config_file"
    fi
    
    echo ""
done

print_header "Testing Complete"
echo "All example configurations have been tested."
