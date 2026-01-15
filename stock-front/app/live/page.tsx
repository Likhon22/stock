"use client";

import { useStockWebSocket } from "@/lib/hooks/useStockWebSocket"; // Adjust import path
import StockCard from "@/components/StockCard";
import { Loader2, Wifi, WifiOff } from "lucide-react";
import { motion } from "framer-motion";

export default function LivePricingPage() {
    // 1. Use the Hook
    const { stocks, isConnected } = useStockWebSocket("ws://localhost:8082/ws");

    return (
        <div className="space-y-8">
            {/* Header with Live Status Indicator */}
            <div className="flex items-center justify-between border-b border-white/10 pb-6">
                <div>
                    <h1 className="text-4xl font-bold tracking-tight text-white mb-2">Live Market</h1>
                    <p className="text-gray-400">Real-time streaming via WebSocket</p>
                </div>

                {/* Connection Badge */}
                <div className={`flex items-center gap-2 px-4 py-2 rounded-full border ${isConnected
                        ? "bg-green-500/10 border-green-500/20 text-green-500"
                        : "bg-red-500/10 border-red-500/20 text-red-500"
                    }`}>
                    {isConnected ? <Wifi className="w-5 h-5" /> : <WifiOff className="w-5 h-5" />}
                    <span className="font-bold uppercase text-sm tracking-wider">
                        {isConnected ? "Connected" : "Disconnected"}
                    </span>
                </div>
            </div>

            {/* Loading State (Only if we have no data yet) */}
            {!isConnected && stocks.length === 0 && (
                <div className="flex flex-col items-center justify-center h-[50vh] text-gray-500 gap-4">
                    <Loader2 className="w-10 h-10 animate-spin text-blue-500" />
                    <p>Establishing secure connection to Wall Street...</p>
                </div>
            )}

            {/* Grid of Live Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {stocks.map((stock) => (
                    // We don't need entry animations here because they would flash every update
                    // Just render the component which handles its own internal numbers
                    <div key={stock.symbol} className="transition-all duration-300">
                        <StockCard stock={stock} />
                    </div>
                ))}
            </div>
        </div>
    );
}