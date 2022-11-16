#!/bin/sh

/app/report-service -config /app/config.yaml | tee /app/report-service.log
