#!/bin/bash

# Create a new user in AWS IAM
read -p "Enter a user name: " user_name
aws iam create-user --user-name "$user_name" --output json
echo "User $user_name created successfully."

# Generate access key and secret key for the user
credentials=$(aws iam create-access-key --user-name "$user_name" --query 'AccessKey.[AccessKeyId,SecretAccessKey]' --output text)
access_key_id=$(echo "$credentials" | cut -f1)
secret_access_key=$(echo "$credentials" | cut -f2)
echo "Access key ID: $access_key_id"
echo "Secret access key: $secret_access_key"

# Attach existing policy to the user
read -p "Enter the policy ARN: " policy_arn
aws iam attach-user-policy --user-name "$user_name" --policy-arn "$policy_arn"
echo "Policy $policy_arn attached to user $user_name."

# Save the credentials in a file
home_dir=$(eval echo ~$USER)
credentials_file="$home_dir/.aws/credentials"
config_file="$home_dir/.aws/config"
read -p "Enter the region: " region
read -p "Enter the output format: " output_format

cat > "$credentials_file" << EOF
[default]
aws_access_key_id = $access_key_id
aws_secret_access_key = $secret_access_key
EOF

cat > "$config_file" << EOF
[default]
region = $region
output = $output_format
EOF

echo "Credentials and config files saved in $home_dir/.aws folder."