import axios from 'axios';

// Ensure this matches your backend URL
const API_BASE_URL = 'http://localhost:3000';

export const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const getPrices = async () => {
    const response = await api.get('/prices');
    return response.data;
};

export const getPrice = async (symbol: string) => {
    const response = await api.get(`/price/${symbol}`);
    return response.data;
};

export const getHistory = async (symbol: string, limit: number = 20) => {
    const response = await api.get(`/price/history/${symbol}/${limit}`);
    return response.data;
};
