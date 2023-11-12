const fs = require('fs');
const path = require('path');

// Function to convert camelCase to snake_case
function camelToSnake(name) {
  return name.replace(/\.?([A-Z]+)/g, (x, y) => "_" + y.toLowerCase()).replace(/^_/, "");
}

// Function to process the content of protobuf files and replace camelCase with snake_case
function processProtobuf(content) {
  // Use a regex to find lines with fields
  const fieldLines = content.match(/^(?!message|})(.*=.*)$/gm);

  if (fieldLines) {
    fieldLines.forEach(line => {
      // Extract the field name (the word before the '=' sign that isn't a type)
      const fieldNameMatch = line.match(/\b(\w+)\s*=/);
      if (fieldNameMatch) {
        const fieldName = fieldNameMatch[1];
        // Check if the field name is camelCase
        if (/[a-z][A-Z]/.test(fieldName)) {
          // Convert the field name to snake_case
          const snakeCaseName = camelToSnake(fieldName);
          // Replace the field name in the line
          content = content.replace(line, line.replace(fieldName, snakeCaseName));
        }
      }
    });
  }
  return content;
}

// Function to recursively read all .proto files from a directory and its subdirectories
function readProtoFiles(dir) {
  const entries = fs.readdirSync(dir, { withFileTypes: true });

  // Array to hold all the .proto files
  let protoFiles = [];

  for (let entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      // If entry is a directory, recursively search it for .proto files
      protoFiles = protoFiles.concat(readProtoFiles(fullPath));
    } else if (path.extname(entry.name) === '.proto') {
      // If entry is a .proto file, add it to the array
      protoFiles.push(fullPath);
    }
  }

  return protoFiles;
}

// Directory of your protobuf files
const protoDir = 'proto/elys';

// Recursively read all .proto files
const protoFiles = readProtoFiles(protoDir);

// Process each .proto file
protoFiles.forEach(filePath => {
  const content = fs.readFileSync(filePath, { encoding: 'utf8' });
  const processedContent = processProtobuf(content);
  fs.writeFileSync(filePath, processedContent);
  console.log(`Processed: ${filePath}`);
});

console.log('All .proto files have been processed.');
