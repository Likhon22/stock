export interface Stock {
    symbol: string;
    name?: string; // API might not return this, using symbol as name mostly
    price: number;
}

export interface StockHistoryItem {
    timestamp: string; // or Date
    price: number;
}

export type StockHistory = StockHistoryItem[];

export interface ApiResponse<T> {
    message: string;
    data: T;
}
