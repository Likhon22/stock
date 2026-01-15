import type { Metadata } from "next";
import { Inter } from "next/font/google"; // Using Inter as standard
import Navbar from "@/components/Navbar";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "StockMaster | Real-time Analytics",
  description: "Advanced stock tracking dashboard",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.className} min-h-screen bg-background text-foreground selection:bg-blue-500/30`}>
        <Navbar />
        <main className="pt-24 pb-12 px-6 max-w-7xl mx-auto">
          {children}
        </main>
      </body>
    </html>
  );
}
