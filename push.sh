#!/usr/bin/env bash
if [[ -z "$1" ]]
then
    echo "Provide commit message"
    exit 1
fi

git add .
git commit -m $1
git push heroku master