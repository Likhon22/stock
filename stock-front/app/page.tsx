"use client";

import { useEffect, useState } from "react";
import { getPrices } from "@/lib/api";
import StockCard from "@/components/StockCard";
import { Stock } from "@/types";
import { motion } from "framer-motion";
import { Loader2 } from "lucide-react";

export default function Dashboard() {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const response = await getPrices();
        // API response wrapper check
        let data: Stock[] = [];

        // Check if response.data is the payload Array
        const payload = response.data || response;

        if (Array.isArray(payload)) {
          console.log("Valid Array Payload:", payload); // Debugging
          data = payload;
        } else {
          console.warn("Unexpected API response structure", payload);
        }
        setStocks(data);
      } catch (err) {
        console.error("Failed to fetch stocks", err);
        setError("Failed to load stock prices. Ensure the backend is running.");
      } finally {
        setLoading(false);
      }
    }

    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
      </div>
    );
  }
  console.log(stocks);


  if (error) {
    return (
      <div className="flex h-[50vh] flex-col items-center justify-center gap-4">
        <div className="p-4 rounded-lg bg-red-500/10 border border-red-500/20 text-red-500">
          {error}
        </div>
        <button
          onClick={() => window.location.reload()}
          className="text-sm underline hover:text-white transition-colors"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight text-white mb-2">Market Overview</h1>
          <p className="text-gray-400">Real-time stock prices and trends</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {stocks.map((stock, index) => (
          <motion.div
            key={stock.symbol}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.05 }}
          >
            <StockCard stock={stock} />
          </motion.div>
        ))}
      </div>

      {stocks.length === 0 && (
        <div className="text-center py-20 text-gray-500">
          No stocks found.
        </div>
      )}
    </div>
  );
}
