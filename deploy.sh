#!/usr/bin/env bash

gcloud builds submit --tag gcr.io/telvanni-platform/discord-bot
gcloud run deploy discord-bot --image gcr.io/telvanni-platform/discord-bot --platform managed --allow-unauthenticated \
  --region us-central1 --set-env-vars="PUBLIC_KEY=dc2a2ef24d22c445bd5a81bab30219e7f1ebbaa8035513457cac4b145b32cdc3"