#!/bin/sh

#Script to run every position in a suite (an EPD file)
#Usage: ./runSuite.sh suite.epd circuit://<go circuit address>

if [ -z $CIRCUIT_ADDR ]; then
    CIRCUIT_ADDR=$2
fi
    
while IFS='' read -r line || [[ -n "$line" ]]; do
    goChess $CIRCUIT_ADDR "$line"
    sleep 1
done < "$1"
