import { ExportToCsv } from 'export-to-csv';

const options = { 
    fieldSeparator: ',',
    quoteStrings: '"',
    decimalSeparator: '.',
    showLabels: true, 
    showTitle: true,
    title: 'Non-linear eden allocation',
    useTextFile: false,
    useBom: true,
    useKeysAsHeaders: true,
  };

export const GenerateCSV = (data: any) => {
    const csvExporter = new ExportToCsv(options);
    const fs = require('fs')
    const csvData = csvExporter.generateCsv(data, true)
    fs.writeFileSync('data.csv',csvData)
}