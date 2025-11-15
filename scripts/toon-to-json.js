#!/usr/bin/env node

import { decode } from '@toon-format/toon';
import { readFileSync } from 'fs';

// Read TOON from stdin or file
let input = '';

if (process.argv[2]) {
  // Read from file
  input = readFileSync(process.argv[2], 'utf8');
} else {
  // Read from stdin (not implemented for simplicity in first version)
  console.error('Please provide a file path');
  process.exit(1);
}

try {
  // Decode TOON to JavaScript object
  const data = decode(input);

  // Output as JSON
  console.log(JSON.stringify(data, null, 2));
} catch (error) {
  console.error('Error parsing TOON:', error.message);
  process.exit(1);
}
