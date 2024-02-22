#!/bin/bash
echo "Script executed: $0 $@" > /root/script_execution.log

input_file=/root/deploy.log
# Define the search string
taiko_l1_address_search_string="^\s*-\s*taikoL1Addr\s*:\s*"
taiko_l1_token_address_search_string="^\s*>\s*taiko_token\s*@\s*"

debug_mode=false

usage() {
  echo "Usage: $0 [-d] TAIKO_L1_ADDRESS|TAIKO_L1_TOKEN_ADDRESS" >&2
  exit 1
}

while getopts ":d" opt; do
  case $opt in
    d)
      debug_mode=true
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      usage
      ;;
  esac
done

shift $((OPTIND - 1))

if [[ $# -ne 1 ]]; then
  usage
fi

# Process command-line options
case "$1" in
TAIKO_L1_ADDRESS)
  cmd="TAIKO_L1_ADDRESS=\$(strings \"$input_file\" | grep -m 1 -e \"$taiko_l1_address_search_string\" | awk '{print \$NF}')"
  if [[ $debug_mode == true ]]; then
    echo "$cmd"
  fi
  eval "$cmd"
  if [[ -n "$TAIKO_L1_ADDRESS" ]]; then
    echo $TAIKO_L1_ADDRESS
    exit 0
  else
    echo "TAIKO_L1_ADDRESS not found"
    exit 1
  fi
  ;;
TAIKO_L1_TOKEN_ADDRESS)
  cmd="TAIKO_L1_TOKEN_ADDRESS=\$(strings \"$input_file\" | grep -A1 -e \"$taiko_l1_token_address_search_string\" | grep -m 1 addr | awk '{print \$NF}')"
  if [[ $debug_mode == true ]]; then
      echo "$cmd"
  fi
  eval "$cmd"
  if [[ -n "$TAIKO_L1_TOKEN_ADDRESS" ]]; then
    echo TAIKO_L1_TOKEN_ADDRESS
    exit 0
  else
    echo "TAIKO_L1_TOKEN_ADDRESS not found"
    exit 1
  fi
  ;;
*)
  echo "Unsupported address: $1" >&2
  exit 1
  ;;
esac
