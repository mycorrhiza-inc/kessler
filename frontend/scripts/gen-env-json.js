#!/usr/bin/env node
const fs = require('fs');
const path = require('path');

const PREFIX = 'NEXT_PUBLIC_';
const envVars = Object.entries(process.env)
  .filter(([key]) => key.startsWith(PREFIX))
  .reduce((acc, [key, value]) => {
    acc[key] = value;
    return acc;
  }, {});

const outputPath = path.resolve(__dirname, '../public/env.json');
fs.writeFileSync(outputPath, JSON.stringify(envVars, null, 2));
console.log(`Wrote public/env.json with ${Object.keys(envVars).length} keys at ${outputPath}`);
