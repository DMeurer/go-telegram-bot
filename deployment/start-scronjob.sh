crontab -u raspi -e # Add the following line to the crontab
* * * * * /home/raspi/go-telegram-bot/deployment/start-scronjob.sh
