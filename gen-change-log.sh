#!/bin/bash

# Get all tags sorted by version
tags=$(git tag --sort=version:refname)

# Check if there are at least two tags
if [ $(echo "$tags" | wc -l) -lt 2 ]; then
  echo "Not enough tags to compare."
  exit 1
fi

# Loop through the tags and show commit logs between each pair of consecutive tags
prev_tag=""
for tag in $tags; do
  if [ -n "$prev_tag" ]; then
    echo "Showing commit logs from $prev_tag to $tag:"
    git log --oneline --graph --decorate $prev_tag..$tag
    echo "---------------------------------------------------"
  fi
  prev_tag=$tag
done
