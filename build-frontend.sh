#!/bin/bash

export VITE_APP_VERSION=$1

cd server/frontend
npm install
npm run build
