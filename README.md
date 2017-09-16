# Go Remote IP Checker

A simple Go program to check your remote IP and e-mail you through your gmail account if it changes.

First, follow these steps to get your client_secrect.json

1. Use this wizard to create or select a project in the Google Developers Console and automatically turn on the API. Click Continue, then Go to credentials.
2. On the Add credentials to your project page, click the Cancel button.
3. At the top of the page, select the OAuth consent screen tab. Select an Email address, enter a Product name if not already set, and click the Save button.
4. Select the Credentials tab, click the Create credentials button and select OAuth client ID.
Select the application type Other, enter the name "Gmail API Quickstart", and click the Create button.
5. Click OK to dismiss the resulting dialog.
6. Click the file_download (Download JSON) button to the right of the client ID.
7. Move this file to your working directory and rename it client_secret.json.

Put your gmail information into config.json:

```json
{
  "From": "you@yourgmail.com",
  "To": "you@yourgmail.com"
}
```

Run this program in a scheduler (cron, etc) and get e-mail updates if/when your remote IP changes.
