#!/usr/bin/env bash

gcloud functions deploy discord-bot --runtime nodejs14 --trigger-http --allow-unauthenticated --entry-point=bot