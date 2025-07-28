#!/usr/bin/env python3

def parse_balances(balance_text):
    """Parse balance text and return a dictionary of denom -> amount."""
    balances = {}
    lines = balance_text.strip().split('\n')
    
    current_amount = None
    for line in lines:
        line = line.strip()
        if line.startswith('- amount:') or line.startswith('amount:'):
            # Extract amount
            current_amount = int(line.split('"')[1])
        elif line.startswith('denom:'):
            # Extract denom
            denom = line.split('denom: ')[1].strip()
            balances[denom] = current_amount
    
    return balances

def format_number(num):
    """Format large numbers with commas."""
    return f"{num:,}"

def main():
    # First observation
    first_observation = """- amount: "485558"
  denom: ibc/0E293A7622DC9A6439DB60E6D234B5AF446962E27CA3AB44D0590603DFF6968E
- amount: "1524142"
  denom: ibc/147B3FF1D005512CCE4089559AF5D0C951F4211A031F15E782E505B85022DF89
- amount: "6910109972376"
  denom: ibc/2A20C613F2BA256BE9FAD791B0743FBC4FE3C998EDF9234D969172F7D42FB8E7
- amount: "86143915"
  denom: ibc/343182F79E6450836403252D1122288D480605885A01426085859B43F5ECD3EF
- amount: "1032815955146700950"
  denom: ibc/421AE2B1D80E42890BF227DBC56C3CADEDA050F21AD1B5DDC1EBE9B9737A4DC6
- amount: "677190"
  denom: ibc/45D6B52CAD911A15BD9C2F5FFDA80E26AFCB05C7CD520070790ABC86D2B24229
- amount: "1029517754220729"
  denom: ibc/622E4D024777290E25ACFB32CE149BB21EF785D78A0A3F297A7AF29A88325AED
- amount: "64100912"
  denom: ibc/646315E3B0461F5FA4C5C8968A88FC45D4D5D04A45B98F1B8294DD82F386DD85
- amount: "863084003251984"
  denom: ibc/68B5E8DA9270FA00245484BB1C07AD75399AD67D54A1344F6E998B5FB69B664F
- amount: "117965145624847"
  denom: ibc/694A6B26A43A2FBECCFFEAC022DEACB39578E54207FDD32005CD976B57B98004
- amount: "6947667"
  denom: ibc/6BFB09FE2464A7681645610F56BBEFF555A00B8AE146339FEB4609BF40FB0F4A
- amount: "8793284"
  denom: ibc/799FDD409719A1122586A629AE8FCA17380351A51C1F47A80A1B8E7F2A491098
- amount: "215862735375812281"
  denom: ibc/8464A63954C0350A26C8588E20719F3A0AC8705E4CA0F7450B60C3F16B2D3421
- amount: "8243363"
  denom: ibc/8BFE59DCD5A7054F0A97CF91F3E3ABCA8C5BA454E548FA512B729D4004584D47
- amount: "193651334361986894"
  denom: ibc/ADF401C952ADD9EE232D52C8303B8BE17FE7953C8D420F20769AF77240BD0C58
- amount: "4730881"
  denom: ibc/B870E6642B6491779D35F326A895780FC2F7409DFD7F639A98505555AEAF345F
- amount: "945266"
  denom: ibc/B88C39AD6C8550716DFD64C3AD28F355633554821249AC9F8BCC21341641CD18
- amount: "15910878"
  denom: ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9
- amount: "6500254"
  denom: ibc/E3459360643C2555C57C7DAB0567FA762B42D5D6D45A76615EA7D99D933AEC04
- amount: "187313800"
  denom: ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349
- amount: "115162249"
  denom: uelys"""

    # Second observation
    second_observation = """- amount: "182"
  denom: ibc/0E293A7622DC9A6439DB60E6D234B5AF446962E27CA3AB44D0590603DFF6968E
- amount: "572"
  denom: ibc/147B3FF1D005512CCE4089559AF5D0C951F4211A031F15E782E505B85022DF89
- amount: "6910109972376"
  denom: ibc/2A20C613F2BA256BE9FAD791B0743FBC4FE3C998EDF9234D969172F7D42FB8E7
- amount: "32304"
  denom: ibc/343182F79E6450836403252D1122288D480605885A01426085859B43F5ECD3EF
- amount: "387305983180013"
  denom: ibc/421AE2B1D80E42890BF227DBC56C3CADEDA050F21AD1B5DDC1EBE9B9737A4DC6
- amount: "254"
  denom: ibc/45D6B52CAD911A15BD9C2F5FFDA80E26AFCB05C7CD520070790ABC86D2B24229
- amount: "1029517754220729"
  denom: ibc/622E4D024777290E25ACFB32CE149BB21EF785D78A0A3F297A7AF29A88325AED
- amount: "24038"
  denom: ibc/646315E3B0461F5FA4C5C8968A88FC45D4D5D04A45B98F1B8294DD82F386DD85
- amount: "323656501219"
  denom: ibc/68B5E8DA9270FA00245484BB1C07AD75399AD67D54A1344F6E998B5FB69B664F
- amount: "117965145624847"
  denom: ibc/694A6B26A43A2FBECCFFEAC022DEACB39578E54207FDD32005CD976B57B98004
- amount: "2605"
  denom: ibc/6BFB09FE2464A7681645610F56BBEFF555A00B8AE146339FEB4609BF40FB0F4A
- amount: "3297"
  denom: ibc/799FDD409719A1122586A629AE8FCA17380351A51C1F47A80A1B8E7F2A491098
- amount: "80948525765930"
  denom: ibc/8464A63954C0350A26C8588E20719F3A0AC8705E4CA0F7450B60C3F16B2D3421
- amount: "3091"
  denom: ibc/8BFE59DCD5A7054F0A97CF91F3E3ABCA8C5BA454E548FA512B729D4004584D47
- amount: "72619250385745"
  denom: ibc/ADF401C952ADD9EE232D52C8303B8BE17FE7953C8D420F20769AF77240BD0C58
- amount: "1774"
  denom: ibc/B870E6642B6491779D35F326A895780FC2F7409DFD7F639A98505555AEAF345F
- amount: "945266"
  denom: ibc/B88C39AD6C8550716DFD64C3AD28F355633554821249AC9F8BCC21341641CD18
- amount: "5967"
  denom: ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9
- amount: "2438"
  denom: ibc/E3459360643C2555C57C7DAB0567FA762B42D5D6D45A76615EA7D99D933AEC04
- amount: "402161"
  denom: ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349
- amount: "11272531150"
  denom: uelys"""

    # Parse both observations
    balances1 = parse_balances(first_observation)
    balances2 = parse_balances(second_observation)
    
    # Calculate differences
    all_denoms = set(balances1.keys()) | set(balances2.keys())
    
    print("BALANCE DIFFERENCE ANALYSIS")
    print("=" * 100)
    print(f"{'Token':<50} {'First Obs':<20} {'Second Obs':<20} {'Difference':<20} {'% Change':<10}")
    print("-" * 100)
    
    increases = []
    decreases = []
    unchanged = []
    
    for denom in sorted(all_denoms):
        amount1 = balances1.get(denom, 0)
        amount2 = balances2.get(denom, 0)
        diff = amount2 - amount1
        
        # Calculate percentage change
        if amount1 != 0:
            pct_change = (diff / amount1) * 100
            pct_str = f"{pct_change:.2f}%"
        else:
            pct_str = "N/A"
        
        # Shorten denom for display
        display_denom = denom[:47] + "..." if len(denom) > 50 else denom
        
        print(f"{display_denom:<50} {format_number(amount1):<20} {format_number(amount2):<20} {format_number(diff):<20} {pct_str:<10}")
        
        if diff > 0:
            increases.append((denom, diff))
        elif diff < 0:
            decreases.append((denom, diff))
        else:
            unchanged.append(denom)
    
    print("\n" + "=" * 100)
    print("SUMMARY:")
    print(f"Tokens with increases: {len(increases)}")
    print(f"Tokens with decreases: {len(decreases)}")
    print(f"Tokens unchanged: {len(unchanged)}")
    
    if increases:
        largest_increase = max(increases, key=lambda x: x[1])
        print(f"\nLargest increase: {largest_increase[0]}")
        print(f"  Change: {format_number(largest_increase[1])}")
    
    if decreases:
        largest_decrease = min(decreases, key=lambda x: x[1])
        print(f"\nLargest decrease: {largest_decrease[0]}")
        print(f"  Change: {format_number(largest_decrease[1])}")
    
    print(f"\nUnchanged tokens: {', '.join(unchanged)}")

if __name__ == "__main__":
    main() 