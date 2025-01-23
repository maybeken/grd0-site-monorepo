#!/bin/bash
aws lambda update-function-code --function-name api-grd0-net --zip-file fileb://build/bootstrap.zip
