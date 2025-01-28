"use client";

import React, {
  createContext,
  useContext,
  useState,
  useCallback,
  useMemo,
} from "react";
import { DrinkLog } from "@/lib/types/api";
import { apiService } from "@/lib/api";
import { format } from "date-fns";

interface DailyStats {
  date: string;
  drinks: DrinkLog[];
  drinkCount: number;
  standardDrinks: number;
}

interface DrinkLogsContextType {
  drinkLogs: DrinkLog[];
  refreshDrinkLogs: () => Promise<void>;
  fetchMoreDrinkLogs: () => Promise<DrinkLog[]>;
  hasMoreLogs: boolean;
  groupedDrinks: Record<string, DrinkLog[]>;
  dailyStats: DailyStats[];
}

const PAGE_SIZE = 10;

const DrinkLogsContext = createContext<DrinkLogsContextType | undefined>(
  undefined
);

export function DrinkLogsProvider({ children }: { children: React.ReactNode }) {
  const [drinkLogs, setDrinkLogs] = useState<DrinkLog[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [hasMoreLogs, setHasMoreLogs] = useState(true);

  // Calculate grouped drinks
  const groupedDrinks = useMemo(() => {
    const groups: Record<string, DrinkLog[]> = {};

    if (!drinkLogs?.length) return groups;

    // Sort all drinks first by date (newest first)
    const sortedDrinks = [...drinkLogs].sort(
      (a, b) =>
        new Date(b.logged_at).getTime() - new Date(a.logged_at).getTime()
    );

    // Group drinks by date
    for (const drink of sortedDrinks) {
      try {
        const drinkDate = new Date(drink.logged_at);
        const dateKey = format(drinkDate, "yyyy-MM-dd");

        if (!groups[dateKey]) {
          groups[dateKey] = [];
        }
        groups[dateKey].push(drink);
      } catch (error) {
        console.error("Error processing drink:", drink, error);
      }
    }

    return groups;
  }, [drinkLogs]);

  // Calculate daily stats
  const dailyStats = useMemo(() => {
    const stats = Object.entries(groupedDrinks).map(([date, drinks]) => {
      const standardDrinks = drinks.reduce((total, drink) => {
        const drinks = Number(drink.standard_drinks) || 0;
        return total + drinks;
      }, 0);

      return {
        date,
        drinks,
        drinkCount: drinks.length,
        standardDrinks,
      };
    });

    return stats.sort(
      (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
    );
  }, [groupedDrinks]);

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
        groupedDrinks,
        dailyStats,
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
