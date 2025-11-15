#!/usr/bin/env node

import { encode } from '@toon-format/toon';

// Read JSON from stdin
let input = '';

if (process.stdin.isTTY) {
  console.error('Please provide JSON input via stdin');
  process.exit(1);
}

const chunks = [];
process.stdin.on('data', (chunk) => chunks.push(chunk));
process.stdin.on('end', () => {
  input = Buffer.concat(chunks).toString('utf8');

  try {
    // Parse JSON
    const data = JSON.parse(input);

    // Encode to TOON
    const toonOutput = encode(data);

    // Output TOON
    console.log(toonOutput);
  } catch (error) {
    console.error('Error converting JSON to TOON:', error.message);
    process.exit(1);
  }
});
