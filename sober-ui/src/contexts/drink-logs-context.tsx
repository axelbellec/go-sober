"use client";

import React, { createContext, useContext, useState, useCallback } from "react";
import { DrinkLog } from "@/lib/types/api";
import { apiService } from "@/lib/api";

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
      const response = await apiService.getDrinkLogs();
      setDrinkLogs(response.drink_logs || []);
    } catch (error) {
      console.error("Failed to fetch drink logs:", error);
      setDrinkLogs([]);
    }
  }, []);

  const addDrinkLog = useCallback((newLog: DrinkLog) => {
    setDrinkLogs((prev) => {
      const currentLogs = Array.isArray(prev) ? prev : [];
      return [newLog, ...currentLogs];
    });
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
