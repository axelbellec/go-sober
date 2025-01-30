import { useState } from "react";
import { formatDistanceToNow } from "date-fns";
import { motion } from "framer-motion";
import { DrinkLog } from "@/lib/types/api";
import { DrinkLogForm } from "@/components/forms/drink-log-form";
import { Button } from "@/components/ui/button";
import { apiService } from "@/lib/api";
import { toast } from "sonner";
import { useDrinkLogs } from "@/contexts/drink-logs-context";
import { useScreenSize } from "@/hooks/use-screen-size";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

export function DrinkLogItem({ drink }: { drink: DrinkLog }) {
  const [isEditing, setIsEditing] = useState(false);
  const { refreshDrinkLogs } = useDrinkLogs();
  const screenSize = useScreenSize();

  const handleReLog = async () => {
    try {
      await apiService.createDrinkLog({
        name: drink.name,
        type: drink.type,
        size_value: drink.size_value,
        size_unit: drink.size_unit,
        abv: drink.abv,
      });
      await refreshDrinkLogs();
      toast.success(
        `Logged ${drink.name} (${drink.size_value}${drink.size_unit}, ${(
          drink.abv * 100
        ).toFixed(2)}% ABV) again`
      );
    } catch (error) {
      console.error("Failed to re-log drink:", error);
      toast.error("Failed to log drink again");
    }
  };

  // Safely parse the date and handle invalid dates
  const loggedDate = new Date(drink.logged_at);
  const timeAgo = !isNaN(loggedDate.getTime())
    ? formatDistanceToNow(loggedDate, { addSuffix: true })
    : "Invalid date";

  const renderEditForm = () => (
    <DrinkLogForm
      initialDrinkLog={drink}
      onCancel={() => setIsEditing(false)}
      onDelete={async () => {
        try {
          await apiService.deleteDrinkLog(drink.id);
          await refreshDrinkLogs();
          toast.success(
            `Deleted ${drink.name} (${drink.size_value}${drink.size_unit}, ${(
              drink.abv * 100
            ).toFixed(2)}% ABV)`
          );
          setIsEditing(false);
        } catch (error) {
          console.error("Error deleting drink log:", error);
        }
      }}
      mode="edit"
    />
  );

  if (isEditing && screenSize.greaterThanOrEqual("md")) {
    return (
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.2 }}
        className="rounded-lg border p-4"
      >
        {renderEditForm()}
      </motion.div>
    );
  }

  return (
    <>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.2 }}
        className="flex flex-col gap-2 rounded-lg border p-4"
      >
        <div
          className="hover:bg-muted/50 p-2 -m-2 rounded-lg cursor-pointer"
          onClick={() => setIsEditing(true)}
          role="button"
          tabIndex={0}
        >
          <div className="flex items-start justify-between gap-4">
            <div className="min-w-0 flex-1">
              <h3 className="font-medium line-clamp-2">{drink.name}</h3>
              <p className="text-sm text-muted-foreground">
                {drink.size_value}
                {drink.size_unit}, {(drink.abv * 100).toFixed(2)}% ABV
              </p>
            </div>
            <div className="flex flex-col items-end gap-1 shrink-0">
              <span className="text-sm text-muted-foreground whitespace-nowrap">
                {timeAgo}
              </span>
              <Button
                variant="outline"
                size="sm"
                className="text-muted-foreground hover:text-foreground hover:bg-muted bg-accent text-accent-foreground h-6 px-2 rounded-md transition-colors duration-200 flex items-center gap-1"
                onClick={(e) => {
                  e.stopPropagation();
                  handleReLog();
                }}
              >
                <span className="opacity-70">+1</span>
              </Button>
            </div>
          </div>
        </div>
      </motion.div>

      <Sheet
        open={isEditing && screenSize.lessThan("md")}
        onOpenChange={setIsEditing}
      >
        <SheetContent
          side="bottom"
          className="h-[80vh] rounded-t-[10px] sm:max-w-none"
        >
          <SheetHeader>
            <SheetTitle>Edit Drink</SheetTitle>
          </SheetHeader>
          <div className="mt-4 overflow-y-auto pb-safe">{renderEditForm()}</div>
        </SheetContent>
      </Sheet>
    </>
  );
}
