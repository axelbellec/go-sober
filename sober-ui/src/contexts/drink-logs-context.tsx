"use client";

import React, { createContext, useContext, useState, useCallback } from "react";
import { DrinkLog } from "@/lib/types/api";
import { apiService } from "@/lib/api";

interface DrinkLogsContextType {
  drinkLogs: DrinkLog[];
  refreshDrinkLogs: () => Promise<void>;
  fetchMoreDrinkLogs: () => Promise<DrinkLog[]>;
  hasMoreLogs: boolean;
}

const PAGE_SIZE = 10;

const DrinkLogsContext = createContext<DrinkLogsContextType | undefined>(
  undefined
);

export function DrinkLogsProvider({ children }: { children: React.ReactNode }) {
  const [drinkLogs, setDrinkLogs] = useState<DrinkLog[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [hasMoreLogs, setHasMoreLogs] = useState(true);

  const refreshDrinkLogs = useCallback(async () => {
    try {
      const response = await apiService.getDrinkLogs({
        page: 1,
        page_size: PAGE_SIZE,
      });
      setDrinkLogs(response.drink_logs);
      setCurrentPage(1);
      setHasMoreLogs(response.drink_logs.length === PAGE_SIZE);
    } catch (error) {
      console.error("Failed to fetch drink logs:", error);
      setHasMoreLogs(false);
    }
  }, []);

  const fetchMoreDrinkLogs = useCallback(async () => {
    try {
      const nextPage = currentPage + 1;
      const response = await apiService.getDrinkLogs({
        page: nextPage,
        page_size: PAGE_SIZE,
      });
      setDrinkLogs((prev) => [...prev, ...response.drink_logs]);
      setCurrentPage(nextPage);
      setHasMoreLogs(response.drink_logs.length === PAGE_SIZE);
      return response.drink_logs;
    } catch (error) {
      console.error("Failed to fetch more drink logs:", error);
      setHasMoreLogs(false);
      return [];
    }
  }, [currentPage]);

  return (
    <DrinkLogsContext.Provider
      value={{
        drinkLogs,
        refreshDrinkLogs,
        fetchMoreDrinkLogs,
        hasMoreLogs,
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
