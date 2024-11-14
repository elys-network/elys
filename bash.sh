#!/bin/bash      

##KEEEPERS MODULES

declare -A allow_coverages
allow_coverages["ACCOUNTEDPOOL"]= ${{vars.ACCOUNTEDPOOL}}

echo "Searching for keeper packages:"
find . -type d -path "*/x/*/keeper" -not -path "*/vendor/*"

module_packages=$(find . -type d -path "*/x/*/keeper" -not -path "*/vendor/*")
for module in $module_packages; do
echo "Checking keeper package: $module"
module_coverage=$(go tool cover -func=coverage.txt | grep "$module" | awk '{sum += $3; count++} END {if (count > 0) print sum / count; else print 0}')
if [ -n "$module_coverage" ]; then
    module_name=$(echo "$module" | sed 's:^.*/\([^/]*\)/keeper$:\1:')
    module_name=$(echo "$module_name" | tr '[:lower:]' '[:upper:]')
    allow_coverage=${allow_coverages[$module_name]}

    echo "La palabra extraída es: $module_name"
    if awk 'BEGIN {exit !('$module_coverage' < '$allow_coverage' )}'; then
    echo "Coverage for $module is below ${allow_coverage}: ${module_coverage}%"
    exit 1
    else
    echo "Coverage for $module is ${module_coverage}%  (meets ${allow_coverage}% requirement)"
    fi
else
    echo "No coverage data found for $module"
fi
done