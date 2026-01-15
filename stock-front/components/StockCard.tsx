import Link from "next/link";
import { ArrowUpRight, ArrowDownRight } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Stock } from "@/types";

interface StockCardProps {
    stock: Stock;
}

export default function StockCard({ stock }: StockCardProps) {
    // Mock trend logic since simple API might not return it
    // In a real app we'd calculate this or get it from API
    const isPositive = Math.random() > 0.5;

    return (
        <Link href={`/stocks/${stock.symbol}`}>
            <Card className="group relative overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-purple-500/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />

                <div className="relative z-10 flex items-start justify-between mb-4">
                    <div>
                        <h3 className="text-2xl font-bold tracking-tight text-white mb-1">
                            {stock.symbol}
                        </h3>
                        <p className="text-sm text-gray-400 font-medium">Stock Name</p>
                    </div>
                    <div className={`
            flex items-center gap-1 px-2 py-1 rounded-full text-xs font-bold
            ${isPositive ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'}
          `}>
                        {isPositive ? <ArrowUpRight className="w-3 h-3" /> : <ArrowDownRight className="w-3 h-3" />}
                        2.4%
                    </div>
                </div>

                <div className="relative z-10">
                    <p className="text-3xl font-bold text-white tracking-tighter">
                        ${stock.price.toFixed(2)}
                    </p>
                </div>
            </Card>
        </Link>
    );
}
