"use client";

import React, { createContext, useContext, useState, useCallback } from "react";
import { DrinkLog } from "@/lib/types/api";
import { fetchWithAuth } from "@/lib/utils";

interface DrinkLogsContextType {
  drinkLogs: DrinkLog[];
  addDrinkLog: (log: DrinkLog) => void;
  refreshDrinkLogs: () => Promise<void>;
}

const DrinkLogsContext = createContext<DrinkLogsContextType | undefined>(
  undefined
);

export function DrinkLogsProvider({ children }: { children: React.ReactNode }) {
  const [drinkLogs, setDrinkLogs] = useState<DrinkLog[]>([]);

  const refreshDrinkLogs = useCallback(async () => {
    try {
      const response = await fetchWithAuth(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/drink-logs`
      );
      const data = await response.json();
      setDrinkLogs(data.drink_logs);
    } catch (error) {
      console.error("Failed to fetch drink logs:", error);
    }
  }, []);

  const addDrinkLog = useCallback((newLog: DrinkLog) => {
    setDrinkLogs((prev) => [newLog, ...prev]);
  }, []);

  return (
    <DrinkLogsContext.Provider
      value={{ drinkLogs, addDrinkLog, refreshDrinkLogs }}
    >
      {children}
    </DrinkLogsContext.Provider>
  );
}

export function useDrinkLogs() {
  const context = useContext(DrinkLogsContext);
  if (context === undefined) {
    throw new Error("useDrinkLogs must be used within a DrinkLogsProvider");
  }
  return context;
}
