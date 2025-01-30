import { useState, useCallback } from "react";

export function useDebounce<T extends (...args: Parameters<T>) => ReturnType<T>>(
    callback: T,
    delay: number
) {
    const [timeoutId, setTimeoutId] = useState<NodeJS.Timeout | null>(null);

    return useCallback(
        (...args: Parameters<T>) => {
            if (timeoutId) {
                clearTimeout(timeoutId);
            }

            setTimeoutId(
                setTimeout(() => {
                    callback(...args);
                    setTimeoutId(null);
                }, delay)
            );
        },
        [callback, delay, timeoutId]
    );
} 