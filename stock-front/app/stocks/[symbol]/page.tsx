"use client";

import { useEffect, useState, use } from "react";
import { getHistory, getPrice } from "@/lib/api";
import { Stock, StockHistory } from "@/types";
import StockChart from "@/components/StockChart";
import { Loader2, ArrowLeft } from "lucide-react";
import Link from "next/link";

interface PageProps {
    params: Promise<{ symbol: string }>;
}

export default function StockDetailPage({ params }: PageProps) {
    // Unwrap params using React.use() or await in async component (Next.js 15+)
    const { symbol } = use(params);

    const [stock, setStock] = useState<Stock | null>(null);
    const [history, setHistory] = useState<StockHistory>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchData() {
            try {
                const [priceData, historyData] = await Promise.all([
                    getPrice(symbol),
                    getHistory(symbol, 50) // Get last 50 points
                ]);

                // Handle API wrapper "data" field
                setStock(priceData.data || priceData);
                setHistory(historyData.data || historyData);
            } catch (err) {
                console.error(err);
            } finally {
                setLoading(false);
            }
        }
        fetchData();
    }, [symbol]);

    if (loading) {
        return (
            <div className="flex h-[50vh] items-center justify-center">
                <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
            </div>
        );
    }

    if (!stock) {
        return (
            <div className="text-center py-20 text-gray-500">
                Stock not found.
            </div>
        );
    }

    return (
        <div className="space-y-8 animate-in fade-in duration-500">
            <Link
                href="/"
                className="inline-flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors"
            >
                <ArrowLeft className="w-4 h-4" />
                Back to Dashboard
            </Link>

            <div className="flex items-end justify-between border-b border-white/10 pb-6">
                <div>
                    <h1 className="text-5xl font-bold text-white mb-2">{stock.symbol}</h1>
                    <p className="text-lg text-gray-400">Real-time Performance</p>
                </div>
                <div className="text-right">
                    <p className="text-4xl font-bold text-blue-500">${stock.price.toFixed(2)}</p>
                    <p className="text-sm text-gray-400">Current Price</p>
                </div>
            </div>

            <div className="glass-card rounded-xl p-6">
                <h2 className="text-xl font-semibold text-white mb-6">Price History</h2>
                <StockChart data={history} />
            </div>
        </div>
    );
}
