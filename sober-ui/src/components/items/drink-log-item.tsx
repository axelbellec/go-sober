import { useState } from "react";
import { formatDistanceToNow } from "date-fns";
import { DrinkLog } from "@/lib/types/api";
import { DrinkLogForm } from "@/components/forms/drink-log-form";
import { Button } from "@/components/ui/button";
import { apiService } from "@/lib/api";
import { toast } from "sonner";
import { useDrinkLogs } from "@/contexts/drink-logs-context";

export function DrinkLogItem({ drink }: { drink: DrinkLog }) {
  const [isEditing, setIsEditing] = useState(false);
  const { refreshDrinkLogs } = useDrinkLogs();

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
      toast.success("Drink logged again");
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

  if (isEditing) {
    return (
      <div className="rounded-lg border p-4">
        <DrinkLogForm
          initialDrinkLog={drink}
          onCancel={() => setIsEditing(false)}
          onDelete={async () => {
            try {
              await apiService.deleteDrinkLog(drink.id);
              await refreshDrinkLogs();
              toast.success("Drink log deleted");
            } catch (error) {
              console.error("Error deleting drink log:", error);
            }
          }}
          mode="edit"
        />
      </div>
    );
  }

  return (
    <div className="rounded-lg border p-4 group relative">
      <div
        className="hover:bg-muted/50 p-2 -m-2 rounded-lg cursor-pointer"
        onClick={() => setIsEditing(true)}
        role="button"
        tabIndex={0}
      >
        <div className="flex justify-between">
          <div>
            <h3 className="font-medium">{drink.name}</h3>
            <p className="text-sm text-muted-foreground">
              {drink.size_value}
              {drink.size_unit}, {(drink.abv * 100).toFixed(2)}% ABV
            </p>
          </div>
          <div className="flex flex-col items-end gap-1">
            <span className="text-sm text-muted-foreground">{timeAgo}</span>
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
    </div>
  );
}
