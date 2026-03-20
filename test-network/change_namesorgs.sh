#!/bin/bash

echo "=========================================================="
echo "TEST: Content replacement respecting capitalization"
echo "=========================================================="

# Search through configuration files
find . -type f \( -name "*.yaml" -o -name "*.sh" -o -name "*.json" \) | while read -r FILE; do
    
    if grep -qiE "orgregistrocivil|orgregistropolicial|orgregistropropiedad" "$FILE"; then
        echo "File: $FILE"
        
        grep -iE "orgregistrocivil|orgregistropolicial|orgregistropropiedad" "$FILE" | while read -r LINE; do
            # 1. First replace capitalized versions (Org)
            # 2. Then lowercase versions (org)
            NEW_LINE=$(echo "$LINE" | \
                sed 's/Orgregistrocivil/Orgregistrocivil/g' | sed 's/orgregistrocivil/orgregistrocivil/g' | \
                sed 's/Orgregistropolicial/Orgregistropolicial/g'           | sed 's/orgregistropolicial/orgregistropolicial/g' | \
                sed 's/Orgregistropropiedad/Orgregistropropiedad/g'   | sed 's/orgregistropropiedad/orgregistropropiedad/g')
            
            echo "  [CURRENT] $LINE"
            echo "  [NEW    ] $NEW_LINE"
            echo "  ------------------------------------------------"
        done
    fi
done
