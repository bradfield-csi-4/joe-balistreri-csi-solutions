#!/bin/bash

if [ -e "$1" ]; then
  echo "$1 exists as a file"
fi

if [ -d "$1" ]; then
  echo "$1 exists as a directory"
fi

if [ -r "$1" ]; then
  echo "$1 exists as a directory"
fi
