import fs from 'fs';
import path from 'path';
import { unstable_noStore as noStore } from 'next/cache';
import { headers } from 'next/headers';

export default function EnvVariablesScript() {
  // Ensure this component always runs on server at request time
  noStore();

  const nonce = headers().get('x-nonce') || '';
  // Read the generated env.json file from public directory
  const jsonPath = path.join(process.cwd(), 'public', 'env.json');
  let envJson = '{}';
  try {
    envJson = fs.readFileSync(jsonPath, 'utf-8');
  } catch (err) {
    console.error('[EnvVariablesScript] Failed to read env.json:', err);
  }

  return (
    <script
      id="env-config"
      nonce={nonce}
      dangerouslySetInnerHTML={{ __html: `window.__ENV__ = ${envJson}` }}
    />
  );
}
