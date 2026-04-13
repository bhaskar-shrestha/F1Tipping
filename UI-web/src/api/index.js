const getApiUrl = () => {
  if (process.env.REACT_APP_API_BASE_URL) {
    return process.env.REACT_APP_API_BASE_URL;
  }
  if (typeof window !== 'undefined' && window.location.hostname !== 'localhost') {
    return '/api';
  }
  return 'http://localhost:8080/api';
};

// Custom error class to preserve HTTP status
class APIError extends Error {
  constructor(status, message, data) {
    super(message);
    this.status = status;
    this.data = data;
  }
}

const API = {
  get: (endpoint) => {
    const url = `${getApiUrl()}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`;
    return fetch(url)
      .then(res => {
        if (!res.ok) {
          return res.text().then(text => {
            const err = new APIError(
              res.status,
              `HTTP ${res.status}: ${res.statusText}`,
              text
            );
            throw err;
          });
        }
        return res.json();
      })
      .catch(err => {
        console.error('API GET error:', url, err);
        throw err;
      });
  },
  post: (endpoint, data) => {
    const url = `${getApiUrl()}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`;
    return fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })
      .then(res => {
        if (!res.ok) {
          return res.text().then(text => {
            const err = new APIError(
              res.status,
              `HTTP ${res.status}: ${res.statusText}`,
              text
            );
            throw err;
          });
        }
        return res.json();
      })
      .catch(err => {
        console.error('API POST error:', url, err);
        throw err;
      });
  },
};

export default API;
