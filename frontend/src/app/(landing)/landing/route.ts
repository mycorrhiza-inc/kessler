import { NextResponse } from 'next/server';
import { readFileSync } from 'fs';
import { join } from 'path';

export function GET() {
  const filePath = join(process.cwd(), 'public', 'landing.html');
  const fileContent = readFileSync(filePath, 'utf8');

  return new NextResponse(fileContent, {
    headers: { 'Content-Type': 'text/html' },
  });
}