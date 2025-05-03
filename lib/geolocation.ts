export async function getStateFromIP(ip: string): Promise<string> {
  try {
    const response = await fetch(`http://ip-api.com/json/${ip}?fields=status,message,regionName`);
    const data = await response.json();

    if (data.status === 'success') {
      return data.regionName;
    }
    return 'New York'; // Default fallback
  } catch (error) {
    console.error('Geolocation lookup failed:', error);
    return 'New York';
  }
}

