# Archived
This repository is archived due to lack of Google API support. [Reference](https://issuetracker.google.com/issues/73122477?pli=1)

# SIL_Mirror_Project
Simple script to create Google mirrors for GitHub and BitBucket repos and configure the mirrors with WebHooks

## Local setup
- Copy example.env to local.env and fill in the values
- Create a file excluded.txt containing a list of repositories to exclude

## Obtaining Google project credentials
- Open https://console.cloud.google.com/welcome
- Create a project (how?)
- Create a service account (how?)
- Go to "IAM & Admin"
- Go to "Service Accounts"
- Open the service account
- Open the Keys tab
- Click Add Key and choose "Create a new key"
- Choose JSON and click Create
- Save the file to this folder
- Set GOOGLE_CREDENTIAL_FILE to the name of the file

## Granting delegated access to the service account
- Go to Google Admin Console
- Open Security - Access and data control - API controls
- Click Manage Domain Wide Delegation (near the bottom of the page)
- In a separate browser tab, go to the Google Cloud Console
- Open IAM & Admin
- Click Service Accounts
- Open the applicable service account
- Copy the Unique ID, which should be a 21-digit number.
- Back on the Google Admin Console, click "Add new"
- In the Client ID input, paste the Unique ID copied from the Google Cloud Console
- Paste these scopes in the next input field, one at a time:
  - `https://www.googleapis.com/auth/cloud-platform`
  - `https://www.googleapis.com/auth/source.full_control`
- Click Authorize
