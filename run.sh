#!/bin/bash
if go build ; then
    ./GoCamera
else
    echo "BUILD ERROR"
fi