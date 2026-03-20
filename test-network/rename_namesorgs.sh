#!/bin/bash

echo "=========================================================="
echo "Updating text in files"
echo "=========================================================="

# Buscamos en los archivos de configuración
find . -type f \( -name "*.yaml" -o -name "*.sh" -o -name "*.json" \) | while read -r FILE; do
    
    # Verificamos si el archivo contiene lo que buscamos
    if grep -qiE "orgregistrocivil|orgregistropolicial|orgregistropropiedad" "$FILE"; then
        echo "Modifying: $FILE"
        
        # Aplicamos los cambios directamente (-i) manteniendo la capitalización
        sed -i 's/Orgregistrocivil/Orgregistrocivil/g; s/orgregistrocivil/orgregistrocivil/g' "$FILE"
        sed -i 's/Orgregistropolicial/Orgregistropolicial/g; s/orgregistropolicial/orgregistropolicial/g' "$FILE"
        sed -i 's/Orgregistropropiedad/Orgregistropropiedad/g; s/orgregistropropiedad/orgregistropropiedad/g' "$FILE"
    fi
done

echo "=========================================================="
echo "Content updated"
echo "=========================================================="
