"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { TrendingUp, Activity, BarChart3 } from "lucide-react";
import { cn } from "@/lib/utils";
import { motion } from "framer-motion";

export default function Navbar() {
    const pathname = usePathname();

    const links = [
        { href: "/", label: "Dashboard", Icon: BarChart3 },
        { href: "/live", label: "Live Pricing", Icon: Activity },
    ];

    return (
        <nav className="fixed top-0 left-0 right-0 z-50 px-6 py-4 glass border-b border-white/5">
            <div className="max-w-7xl mx-auto flex items-center justify-between">
                <Link href="/" className="flex items-center gap-2 group">
                    <div className="p-2 rounded-lg bg-blue-600 group-hover:bg-blue-500 transition-colors">
                        <TrendingUp className="w-6 h-6 text-white" />
                    </div>
                    <span className="text-xl font-bold tracking-tight text-white group-hover:text-blue-400 transition-colors">
                        Stock<span className="text-blue-500">Master</span>
                    </span>
                </Link>

                <div className="flex items-center gap-6">
                    {links.map(({ href, label, Icon }) => {
                        const isActive = pathname === href;
                        return (
                            <Link
                                key={href}
                                href={href}
                                className={cn(
                                    "relative flex items-center gap-2 text-sm font-medium transition-colors",
                                    isActive ? "text-blue-400" : "text-gray-400 hover:text-white"
                                )}
                            >
                                <Icon className="w-4 h-4" />
                                <span>{label}</span>
                                {isActive && (
                                    <motion.div
                                        layoutId="navbar-indicator"
                                        className="absolute -bottom-4 left-0 right-0 h-0.5 bg-blue-500 shadow-[0_0_10px_rgba(59,130,246,0.5)]"
                                    />
                                )}
                            </Link>
                        );
                    })}
                </div>
            </div>
        </nav>
    );
}
