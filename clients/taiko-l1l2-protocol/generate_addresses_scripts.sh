#!/bin/bash
script_dir=$(dirname "$(readlink -f "$0")")

base_script="$script_dir/get_deploy_address.sh"

# Extract the usage information from the base script
usage=$(bash "$base_script" 2>&1)

# Extract the supported parameters from the usage information
# Assuming the parameters are specified as "TAIKO_L1_ADDRESS|TAIKO_L1_TOKEN_ADDRESS"
parameters=$(echo "$usage" | grep -oP "TAIKO_[^ ]+" | tr '|' '\n')

# Generate scripts for each supported parameter
for param in $parameters; do
    script_name="$script_dir/get_deploy_address_${param}.sh"

    # Create the script
    echo "#!/bin/bash" > "$script_name"
    echo "bash $base_script $param" >> "$script_name"

    # Make the script executable
    chmod +x "$script_name"
done

