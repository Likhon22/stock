import { cn } from "@/lib/utils";
import { HTMLAttributes, forwardRef } from "react";

const Card = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
    ({ className, ...props }, ref) => (
        <div
            ref={ref}
            className={cn(
                "glass-card rounded-xl p-6 transition-all duration-300 hover:shadow-blue-500/10",
                className
            )}
            {...props}
        />
    )
);
Card.displayName = "Card";

export { Card };
