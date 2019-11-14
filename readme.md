# Introduction

The sole purpose of this project is for my own amusement. This is a Slack bot written in golang that performs various brilliantly useless functions.

* Slash commands
    * /mark \[your message\] - echos your message except rAndoMLY cApitALizInG characters
    * /score \[word++|word--|top\] - tracks the score for something, anything
* RTM responses
    * any images sent to a private / public channel (excluding DMs) will be saved into the corresponding folder in Google Drive

# Getting Started

1. Clone this repository.

1. Create an environment file `.env` and add your Slack API configs. See the [official Slack guide][1] on create a Slack App and retrieve these values!

   ```
   SLACK_OAUTH_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   SLACK_VERIFICATION_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   SLACK_SIGNING_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   ```

1. Create an `etc` subdirectory and copy the Google Drive API configs including credentials.json and token.json. See the [official Google Drive guide][2] on how to enable the Google Drive API on your account and generate these credential files!

1. Your project directory should now look similarly:

   ```
   pikabot/
     data/
     entry/
     etc/
       credentials.json
       token.json
     .env
     ...
   ```

1. Run `docker-compose build` and `docker-compose up -d` to start it up! Some additional helpful scripts are in `scripts/` if you want to use Let's Encrypt, Nginx, and AWS EC2 together.

# References

This was a really helpful reference when getting started with uploading files to Google Drive.

https://medium.com/@devtud/upload-files-in-google-drive-with-golang-and-google-drive-api-d686fb62f884

[1]: https://api.slack.com/
[2]: https://developers.google.com/drive/api/v3/enable-drive-api
