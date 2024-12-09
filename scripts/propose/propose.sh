#!/bin/bash
arr_poolIds=( $(tail -n +4 data.csv | cut -d ',' -f1) )
arr_poolNames=( $(tail -n +4 data.csv | cut -d ',' -f2) )
arr_multipliers=( $(tail -n +4 data.csv | cut -d ',' -f3) )
arr_nextAllocations=( $(tail -n +4 data.csv | cut -d ',' -f4) )

echo "array of PoolIds  : ${arr_poolIds[@]}"
echo "array of PoolNames   : ${arr_poolNames[@]}"
echo "array of Multipliers : ${arr_multipliers[@]}"
echo "array of Next_Eden_Allocation : ${arr_nextAllocations[@]}"

## now loop through the above arr_poolIds
for id in "${arr_poolIds[@]}"
do
   n=${#id}
   if [ $n -gt 0 ]; then
      poolIds+="$id";
      poolIds+=",";
   fi
done

## now loop through the above arr_multipliers
for m in "${arr_multipliers[@]}"
do
   n=${#m}
   if [ $n -gt 0 ]; then
      multipliers+="$m";
      multipliers+=",";
   fi
done

# remove last comma
poolIds=${poolIds:0:${#poolIds}-1}
# remove last comma
multipliers=${multipliers:0:${#multipliers}-1}

# invalid multipliers
if [ ${#multipliers} -lt 1 ]; then
   exit
fi

echo "!! Submit update-pool-info proposal !!"
elysd tx incentive update-pool-info "$poolIds" "$multipliers" --title="title" --description="description" --deposit="10000000uelys" --from=treasury --chain-id=elysicstestnet-1 --broadcast-mode=block --yes --gas auto