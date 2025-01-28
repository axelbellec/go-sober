"use client";

import { useEffect, useState, useMemo } from "react";
import { useDrinkLogs } from "@/contexts/drink-logs-context";
import { format } from "date-fns";
import { motion, AnimatePresence } from "framer-motion";
import { DrinkLog } from "@/lib/types/api";
import { Button } from "@/components/ui/button";
import { DrinkLogItem } from "@/components/items/drink-log-item";
import { Plus, Info, Beer } from "lucide-react";
import { useScreenSize } from "@/hooks/use-screen-size";
import { Badge } from "@/components/ui/badge";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { DrinkLogForm } from "@/components/forms/drink-log-form";

export function ConsumptionHistoryView() {
  const { drinkLogs, refreshDrinkLogs, fetchMoreDrinkLogs, hasMoreLogs } =
    useDrinkLogs();
  const [isLoading, setIsLoading] = useState(false);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const screenSize = useScreenSize();

  useEffect(() => {
    refreshDrinkLogs();
  }, [refreshDrinkLogs]);

  // Move the grouping logic into useMemo
  const groupedDrinks = useMemo(() => {
    return (drinkLogs ?? []).reduce((groups, drink) => {
      const drinkDate = new Date(drink.logged_at);
      if (isNaN(drinkDate.getTime())) {
        return groups;
      }

      const date = format(drinkDate, "yyyy-MM-dd");
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(drink);
      groups[date].sort(
        (a, b) =>
          new Date(b.logged_at).getTime() - new Date(a.logged_at).getTime()
      );
      return groups;
    }, {} as Record<string, DrinkLog[]>);
  }, [drinkLogs]);

  // Calculate daily stats using useMemo
  const dailyStats = useMemo(() => {
    return Object.entries(groupedDrinks).map(([date, drinks]) => ({
      date,
      drinks,
      drinkCount: drinks.reduce((total) => total + 1, 0),
      standardDrinks: drinks.reduce(
        (total, drink) => total + drink.standard_drinks,
        0
      ),
    }));
  }, [groupedDrinks]);

  const handleLoadMore = async () => {
    setIsLoading(true);
    try {
      await fetchMoreDrinkLogs();
    } catch (error) {
      console.error("Error loading more drinks:", error);
    } finally {
      setIsLoading(false);
    }
  };

  // Ensure drinkLogs is an array before checking length
  const hasLogs = Array.isArray(drinkLogs) && drinkLogs.length > 0;

  return (
    <div className="space-y-8 relative min-h-[300px]">
      {dailyStats
        .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
        .map(({ date, drinks, drinkCount, standardDrinks }) => (
          <div key={date} className="space-y-4">
            <h2 className="text-lg font-semibold flex items-center gap-3">
              {format(new Date(date), "EEEE, MMMM d, yyyy")}
            </h2>
            <div className="flex gap-x-2">
              <Badge
                variant="secondary"
                className="flex items-center gap-x-1 text-muted-foreground font-semibold tracking-tight"
              >
                <Beer className="h-4 w-4" />
                <p className="text-sm">
                  {drinkCount} {drinkCount === 1 ? "drink" : "drinks"}
                </p>
              </Badge>
              <Badge
                variant="secondary"
                className="flex items-center gap-x-1 text-muted-foreground font-semibold tracking-tight"
              >
                <Info className="h-4 w-4" />
                <p className="text-sm">
                  {standardDrinks.toFixed(1)} standard drinks
                </p>
              </Badge>
            </div>
            <AnimatePresence>
              <motion.div
                className="space-y-4"
                initial="hidden"
                animate="visible"
                variants={{
                  visible: {
                    transition: {
                      staggerChildren: 0.1,
                    },
                  },
                }}
              >
                {drinks.map((drink) => (
                  <motion.div
                    key={drink.id}
                    variants={{
                      hidden: { opacity: 0, y: 20 },
                      visible: { opacity: 1, y: 0 },
                    }}
                  >
                    <DrinkLogItem drink={drink} />
                  </motion.div>
                ))}
              </motion.div>
            </AnimatePresence>
          </div>
        ))}

      {hasLogs && (
        <div className="flex justify-center py-4">
          {hasMoreLogs ? (
            <Button
              variant="outline"
              onClick={handleLoadMore}
              disabled={isLoading}
            >
              {isLoading ? "Loading..." : "Load More"}
            </Button>
          ) : (
            <span className="text-muted-foreground">No more drinks logged</span>
          )}
        </div>
      )}

      {!hasLogs && (
        <div className="text-center py-8 text-muted-foreground">
          No drinks logged yet
        </div>
      )}

      {/* Floating Action Button and Sheet */}
      <Sheet open={isSheetOpen} onOpenChange={setIsSheetOpen}>
        <SheetTrigger asChild>
          <Button
            size="icon"
            className="fixed bottom-6 right-6 h-14 w-14 rounded-full shadow-lg"
          >
            <Plus className="h-6 w-6" />
          </Button>
        </SheetTrigger>
        <SheetContent
          side={screenSize.greaterThanOrEqual("lg") ? "right" : "bottom"}
          className={`${
            screenSize.greaterThanOrEqual("lg")
              ? "h-full w-[400px] border-l"
              : "h-[80vh] rounded-t-[10px] sm:max-w-none"
          }`}
        >
          <SheetHeader>
            <SheetTitle>Log a Drink</SheetTitle>
          </SheetHeader>
          <div className="mt-4 overflow-y-auto pb-safe">
            <DrinkLogForm
              onCancel={() => setIsSheetOpen(false)}
              mode="create"
            />
          </div>
        </SheetContent>
      </Sheet>
    </div>
  );
}
