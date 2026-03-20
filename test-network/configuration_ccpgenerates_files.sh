#!/bin/bash

echo "=========================================================="
echo "Updating ORG assignments in ccp-generate.sh"
echo "=========================================================="

# 1. addorgregistropropiedad file (Change ORG=3 to ORG=contraloria)
# This path corresponds to the Property Registry infrastructure being deployed.
FILE1="$HOME/fabric-samples/test-network/addorgregistropropiedad/ccp-generate.sh"
if [ -f "$FILE1" ]; then
    sed -i 's/ORG=3/ORG=contraloria/g' "$FILE1"
    echo "✔ Updated ORG=contraloria in: $FILE1"
else
    echo "✘ Error: Not found $FILE1"
fi

# 2. organizations file (Change ORG=1 to registrocivil and ORG=2 to cne)
# These changes reflect the entities in the public records registration system.
FILE2="$HOME/fabric-samples/test-network/organizations/ccp-generate.sh"
if [ -f "$FILE2" ]; then
    sed -i 's/ORG=1/ORG=registrocivil/g' "$FILE2"
    sed -i 's/ORG=2/ORG=cne/g' "$FILE2"
    echo "✔ Updated ORG=registrocivil and ORG=cne in: $FILE2"
else
    echo "✘ Error: Not found $FILE2"
fi

echo "=========================================================="
echo "Changes successfully implemented!"
echo "=========================================================="
