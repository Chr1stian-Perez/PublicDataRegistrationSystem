#!/bin/bash

echo "=========================================================="
echo "Updating text in files"
echo "=========================================================="


find . -type f \( -name "*.yaml" -o -name "*.sh" -o -name "*.json" \) | while read -r FILE; do
    
    
    if grep -qiE "orgregistrocivil|orgcne|orgcontraloria" "$FILE"; then
        echo "Modifying: $FILE"
        
        
        sed -i 's/Orgregistrocivil/Orgregistrocivil/g; s/orgregistrocivil/orgregistrocivil/g' "$FILE"
        sed -i 's/Orgcne/Orgcne/g; s/orgcne/orgcne/g' "$FILE"
        sed -i 's/Orgcontraloria/Orgcontraloria/g; s/orgcontraloria/orgcontraloria/g' "$FILE"
    fi
done

echo "=========================================================="
echo "Content updated"
echo "=========================================================="