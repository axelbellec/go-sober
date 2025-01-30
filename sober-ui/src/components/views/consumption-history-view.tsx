"use client";

import { useEffect, useState } from "react";
import { useDrinkLogs } from "@/contexts/drink-logs-context";
import { format } from "date-fns";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { DrinkLogItem } from "@/components/items/drink-log-item";
import { Plus, Info, Beer } from "lucide-react";
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
  const {
    drinkLogs,
    refreshDrinkLogs,
    fetchMoreDrinkLogs,
    hasMoreLogs,
    dailyStats,
  } = useDrinkLogs();
  const [isLoading, setIsLoading] = useState(false);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const screenSize = useScreenSize();

  useEffect(() => {
    refreshDrinkLogs();
  }, [refreshDrinkLogs]);

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
      {/* Desktop Add Button */}
      <div className="hidden md:flex justify-end mb-4">
        <Sheet open={isSheetOpen} onOpenChange={setIsSheetOpen}>
          <SheetTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-5 w-5" />
              Add Drink
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

      {/* Mobile Floating Action Button and Sheet */}
      <Sheet open={isSheetOpen} onOpenChange={setIsSheetOpen}>
        <SheetTrigger asChild>
          <Button
            size="icon"
            className="md:hidden fixed bottom-6 right-6 h-14 w-14 rounded-full shadow-lg"
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

      {dailyStats.map(({ date, drinks, drinkCount, standardDrinks }) => (
        <div key={date} className="space-y-4">
          <h2 className="text-lg font-semibold flex items-center gap-3">
            {format(new Date(date), "EEEE, MMMM d, yyyy")}
          </h2>
          <div className="flex gap-x-2 justify-start">
            <Sheet>
              <SheetTrigger asChild>
                <Badge
                  variant="secondary"
                  className="flex items-center gap-x-1 text-muted-foreground font-semibold tracking-tight cursor-pointer"
                >
                  <Beer className="h-4 w-4" />
                  <p className="text-sm">
                    {drinkCount} {drinkCount === 1 ? "drink" : "drinks"}{" "}
                    consumed
                  </p>
                </Badge>
              </SheetTrigger>
              <SheetTrigger asChild>
                <Badge
                  variant="secondary"
                  className="flex items-center gap-x-1 text-muted-foreground font-semibold tracking-tight cursor-pointer"
                >
                  <Info className="h-4 w-4" />
                  <p className="text-sm">
                    {standardDrinks.toFixed(1)} standard drinks
                  </p>
                </Badge>
              </SheetTrigger>
              <SheetContent
                side={screenSize.greaterThanOrEqual("lg") ? "right" : "bottom"}
                className={`${
                  screenSize.greaterThanOrEqual("lg")
                    ? "w-[400px] border-l"
                    : "h-[80vh] rounded-t-[10px] sm:max-w-none"
                }`}
              >
                <SheetHeader>
                  <SheetTitle>Understanding Your Drinks</SheetTitle>
                </SheetHeader>
                <div className="mt-6 space-y-6">
                  <div className="space-y-4">
                    <h3 className="font-medium">Drinks Consumed</h3>
                    <p className="text-muted-foreground">
                      This shows the total number of individual drinks you've
                      had, regardless of their size or alcohol content.
                    </p>
                  </div>
                  <div className="space-y-4">
                    <h3 className="font-medium">Standard Drinks</h3>
                    <p className="text-muted-foreground">
                      A standard drink contains 10g of pure alcohol. This
                      measurement helps track alcohol intake more accurately
                      than counting drinks alone.
                    </p>
                    <div className="space-y-3">
                      <p className="font-medium">Common bar servings:</p>
                      <ul className="space-y-2 pl-4 list-disc">
                        <li>🍺 Half pint (250ml) of beer = 1 standard drink</li>
                        <li>🍷 Glass of wine (125ml) = 1 standard drink</li>
                        <li>🍺 Pint (500ml) of beer = 2 standard drinks</li>
                        <li>🥃 Shot of spirits (30ml) = 1 standard drink</li>
                        <li>
                          🍸 Cocktail (60ml) spirit-based = 2 standard drinks
                        </li>
                      </ul>
                    </div>
                  </div>
                </div>
              </SheetContent>
            </Sheet>
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
    </div>
  );
}
