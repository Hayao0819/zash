#!/bin/sh

# Test script
echo "Hello from script"  
ls /home > /dev/null &
echo "Script finished"
exit 0
