"use client";

import React, { createContext, useContext, useState, useCallback } from "react";
import { DrinkLog } from "@/lib/types/api";
import { apiService } from "@/lib/api";

interface DrinkLogsContextType {
  drinkLogs: DrinkLog[];
  refreshDrinkLogs: () => Promise<void>;
  fetchMoreDrinkLogs: () => Promise<DrinkLog[]>;
}

const DrinkLogsContext = createContext<DrinkLogsContextType | undefined>(
  undefined
);

export function DrinkLogsProvider({ children }: { children: React.ReactNode }) {
  const [drinkLogs, setDrinkLogs] = useState<DrinkLog[]>([]);
  const [currentPage, setCurrentPage] = useState(1);

  const refreshDrinkLogs = useCallback(async () => {
    try {
      const response = await apiService.getDrinkLogs({ page: 1 });
      setDrinkLogs(response.drink_logs);
      setCurrentPage(1);
    } catch (error) {
      console.error("Failed to fetch drink logs:", error);
    }
  }, []);

  const fetchMoreDrinkLogs = useCallback(async () => {
    try {
      const nextPage = currentPage + 1;
      const response = await apiService.getDrinkLogs({ page: nextPage });
      setDrinkLogs((prev) => [...prev, ...response.drink_logs]);
      setCurrentPage(nextPage);
      return response.drink_logs;
    } catch (error) {
      console.error("Failed to fetch more drink logs:", error);
      return [];
    }
  }, [currentPage]);

  return (
    <DrinkLogsContext.Provider
      value={{
        drinkLogs,
        refreshDrinkLogs,
        fetchMoreDrinkLogs,
      }}
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
