#!/bin/sh

echo "Testing while loop"
count=0
while [ $count -lt 3 ]; do
    echo "Count: $count"
    count=$((count + 1))
done
echo "Loop finished"
