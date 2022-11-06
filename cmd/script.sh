mkdir /home/igor/vs_project/birdhouse/executable/
go build -o /home/igor/vs_project/birdhouse/executable/
export APP_GUID= ${{ secret.APP_GUID }}
export BASE_PATH= ${{ secret.BASE_PATH }}
export PORT= ${{ secret.PORT }}
export PUBLIC_KEY= ${{ secret.PUBLIC_KEY }}
export TELEGRAM_BOT_TOKEN= ${{ secret.TELEGRAM_BOT_TOKEN }}
export TOKEN_TIME_TO_LIVE= ${{ secret.TOKEN_TIME_TO_LIVE }}
export SEED=${{ secret.SEED }}

cd /home/igor/vs_project/birdhouse/executable/
./cmd
