import { useEffect, useRef, useState } from "react";
import { Stock } from "@/types";

export function useStockWebSocket(url: string) {
    // 1. Store the connection in a Ref so it persists without triggering re-renders
    // Why Ref? Because we don't want the component to re-render just because the socket object changed.
    const socketRef = useRef<WebSocket | null>(null);

    // 2. State for the UI
    const [stocks, setStocks] = useState<Stock[]>([]);
    const [isConnected, setIsConnected] = useState(false);

    useEffect(() => {
        // 3. Connect
        console.log("🔌 Connecting to WebSocket:", url);
        const socket = new WebSocket(url);
        socketRef.current = socket;

        // 4. Lifecycle Handlers
        socket.onopen = () => {
            console.log("✅ WebSocket Connected");
            setIsConnected(true);
        };

        socket.onclose = () => {
            console.log("❌ WebSocket Disconnected");
            setIsConnected(false);
        };

        socket.onerror = (error) => {
            console.error("⚠️ WebSocket Error:", error);
        };

        // 5. The Core Logic: Handling Messages
        socket.onmessage = (event) => {
            try {
                // Parse the incoming single stock update
                const newStock: Stock = JSON.parse(event.data);

                // FUNCTIONAL STATE UPDATE (The most important part)
                // We access 'prevStocks' which guarantees we have the latest state.
                setStocks((prevStocks) => {
                    // Check if this stock is already in our list
                    const index = prevStocks.findIndex(s => s.symbol === newStock.symbol);

                    if (index !== -1) {
                        // UPDATE: Create a NEW array copy (Immutability) with the replaced item
                        const newArray = [...prevStocks];
                        newArray[index] = newStock;
                        return newArray;
                    } else {
                        // INSERT: Append to the list
                        return [...prevStocks, newStock];
                    }
                });
            } catch (err) {
                console.error("Failed to parse WS message", err);
            }
        };

        // 6. Cleanup (The Zombie Killer)
        // This runs when the component unmounts (User leaves page)
        return () => {
            if (socket.readyState === WebSocket.OPEN) {
                socket.close();
            }
        };
    }, [url]); // Only re-run if URL changes

    return { stocks, isConnected };
}