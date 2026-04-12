// Get API base URL from environment variable with smart fallback
const getApiUrl = () => {
  // If explicitly set in env vars, use that
  if (process.env.REACT_APP_API_BASE_URL) {
    return process.env.REACT_APP_API_BASE_URL;
  }
  // In Docker/production (non-localhost), try relative API path
  if (typeof window !== 'undefined' && window.location.hostname !== 'localhost') {
    return '/api';
  }
  // For local development, use localhost:8080
  return 'http://localhost:8080/api';
};

const API_URL = getApiUrl();

const API = {
  get: (endpoint) => fetch(`${API_URL}/${endpoint}`).then(res => res.json()),
  post: (endpoint, data) => fetch(`${API_URL}/${endpoint}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  }).then(res => res.json()),
};

export default API;
