const API_BASE_URL = 'http://localhost:8080/api';

const API = {
  get: (endpoint) => fetch(`${API_BASE_URL}/${endpoint}`).then(res => res.json()),
  post: (endpoint, data) => fetch(`${API_BASE_URL}/${endpoint}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  }).then(res => res.json()),
};

export default API;
