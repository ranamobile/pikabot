# Getting Started

1. Clone this repository.

1. Create an environment file `.env` and add your Slack API configs.

   ```
   SLACK_OAUTH_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   SLACK_VERIFICATION_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   SLACK_SIGNING_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   ```

1. Create an `etc` subdirectory and copy the Google Drive API configs including credentials.json and token.json.

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
