"use client";

import {
    Area,
    AreaChart,
    CartesianGrid,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from "recharts";
import { StockHistory } from "@/types";

interface StockChartProps {
    data: StockHistory;
}

export default function StockChart({ data }: StockChartProps) {
    // Format data for recharts if needed, generally it accepts array of objects
    // Ideally, data assumes { timestamp, price }

    return (
        <div className="h-[400px] w-full mt-6">
            <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={data}>
                    <defs>
                        <linearGradient id="colorPrice" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3} />
                            <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                        </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.05)" />
                    <XAxis
                        dataKey="timestamp"
                        stroke="#525252"
                        tick={{ fill: '#525252', fontSize: 12 }}
                        tickFormatter={(val) => {
                            // Handle both ISO string or Date object
                            const date = new Date(val);
                            // Return HH:MM:SS only
                            return date.toLocaleTimeString('en-US', { hour12: false });
                        }}
                    />
                    <YAxis
                        stroke="#525252"
                        tick={{ fill: '#525252', fontSize: 12 }}
                        domain={['auto', 'auto']}
                        tickFormatter={(val) => `$${val.toFixed(2)}`}
                    />
                    <Tooltip
                        contentStyle={{
                            backgroundColor: '#171717',
                            borderColor: '#262626',
                            borderRadius: '8px',
                            color: '#ededed'
                        }}
                        formatter={(value: number | undefined) => [
                            value !== undefined ? `$${value.toFixed(2)}` : 'N/A',
                            "Price"
                        ]}
                        labelFormatter={(label) => new Date(label).toLocaleTimeString()}
                    />
                    <Area
                        type="monotone"
                        dataKey="price"
                        stroke="#3b82f6"
                        strokeWidth={2}
                        fillOpacity={1}
                        fill="url(#colorPrice)"
                    />
                </AreaChart>
            </ResponsiveContainer>
        </div>
    );
}
