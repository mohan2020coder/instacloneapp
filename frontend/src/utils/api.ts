const API_BASE_URL = 'http://localhost:8082'; // URL of your Go backend

// A utility function to handle API requests
async function apiRequest<T>(
  endpoint: string,
  method: 'GET' | 'POST' = 'GET',
  body?: object
): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method,
    headers: {
      'Content-Type': 'application/json',
    },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
}

export default apiRequest;
