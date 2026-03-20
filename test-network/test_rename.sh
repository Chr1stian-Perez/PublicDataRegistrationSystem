#!/bin/bash

echo "=========================================================="
echo "MODO SIMULACIÓN MEJORADO: Renombrado total de rutas"
echo "=========================================================="

# 1. Primero identificamos todos los elementos que contienen "org"
# 2. Los ordenamos por longitud de ruta (de más largo a más corto) para no romper las rutas
find . -iname "*org*" | sort -r | while read -r FULL_PATH; do
    
    DIR=$(dirname "$FULL_PATH")
    OLD_NAME=$(basename "$FULL_PATH")
    
    # Aplicamos las sustituciones
    NEW_NAME="$OLD_NAME"
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgregistrocivil/orgregistrocivil/gI')
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgregistropolicial/orgregistropolicial/gI')
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgregistropropiedad/orgregistropropiedad/gI')

    if [ "$OLD_NAME" != "$NEW_NAME" ]; then
        echo "CAMBIO DETECTADO:"
        echo "  De: $FULL_PATH"
        echo "  A:  $DIR/$NEW_NAME"
        echo "----------------------------------------------------------"
    fi
done

echo "=========================================================="
echo "Fin de la simulación. Revisa si 'addOrgregistropropiedad' ahora aparece como 'addOrgregistropropiedad'."
