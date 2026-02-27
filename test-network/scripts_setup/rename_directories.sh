#!/bin/bash

echo "=========================================================="
echo "Running directory and file renaming"
echo "=========================================================="


find . -iname "*org*" | sort -r | while read -r FULL_PATH; do
    
    DIR=$(dirname "$FULL_PATH")
    OLD_NAME=$(basename "$FULL_PATH")
    
    
    NEW_NAME="$OLD_NAME"
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgregistrocivil/orgregistrocivil/gI')
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgcne/orgcne/gI')
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgcontraloria/orgcontraloria/gI')

    if [ "$OLD_NAME" != "$NEW_NAME" ]; then
        
        mv "$FULL_PATH" "$DIR/$NEW_NAME"
        echo "LOG: $OLD_NAME changed to $NEW_NAME"
    fi
done

echo "=========================================================="
echo "¡Renaming completed successfully!"
echo "=========================================================="