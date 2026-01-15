import axios from 'axios';

const api = axios.create({
  // baseURL: 'http://127.0.0.1:8080',
  baseURL: 'http://192.168.1.35:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export default api;