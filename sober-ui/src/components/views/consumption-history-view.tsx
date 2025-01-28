"use client";

import { useEffect, useState } from "react";
import { useDrinkLogs } from "@/contexts/drink-logs-context";
import { format } from "date-fns";
import { motion, AnimatePresence } from "framer-motion";
import { DrinkLog } from "@/lib/types/api";
import { Button } from "@/components/ui/button";
import { DrinkLogItem } from "@/components/items/drink-log-item";
import { Plus } from "lucide-react";
import { useScreenSize } from "@/hooks/use-screen-size";
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

  // Group drinks by date
  const groupedDrinks = (drinkLogs ?? []).reduce((groups, drink) => {
    // Safely parse the date and handle invalid dates
    const drinkDate = new Date(drink.logged_at);
    if (isNaN(drinkDate.getTime())) {
      return groups; // Skip invalid dates
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
      {Object.entries(groupedDrinks)
        .sort((a, b) => new Date(b[0]).getTime() - new Date(a[0]).getTime())
        .map(([date, drinks]) => (
          <div key={date} className="space-y-4">
            <h2 className="text-lg font-semibold">
              {format(new Date(date), "EEEE, MMMM d, yyyy")}
            </h2>
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
              : "h-[90vh] rounded-t-[10px] sm:max-w-none"
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
