"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

import { DrinkOption } from "@/lib/types/api";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { apiService } from "@/lib/api";
import { useDrinkLogs } from "@/contexts/drink-logs-context";
import { toast } from "sonner";

const drinkLogSchema = z.object({
  drinkOptionId: z.number().min(1, "Please select a drink"),
  abv: z.number().min(0.01, "ABV must be greater than 0").max(100),
  sizeValue: z.number().min(1, "Size must be greater than 0"),
  sizeUnit: z.string().min(1, "Unit is required"),
});

type DrinkLogFormValues = z.infer<typeof drinkLogSchema>;

export function DrinkLogForm() {
  const [drinkOptions, setDrinkOptions] = useState<DrinkOption[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<DrinkLogFormValues>({
    resolver: zodResolver(drinkLogSchema),
    defaultValues: {
      drinkOptionId: 0,
      abv: 0,
      sizeValue: 0,
      sizeUnit: "cl",
    },
    mode: "onChange",
  });

  const { addDrinkLog, refreshDrinkLogs } = useDrinkLogs();

  useEffect(() => {
    apiService
      .getDrinkOptions()
      .then((data) => {
        const options = Array.isArray(data.drink_options)
          ? data.drink_options
          : [];
        setDrinkOptions(options);
      })
      .catch((error) => {
        console.error("Failed to fetch drink options:", error);
        setDrinkOptions([]);
      });
  }, []);

  const onDrinkSelect = (drinkId: string) => {
    const drink = drinkOptions.find((d) => d.id === parseInt(drinkId));
    if (drink) {
      form.setValue("drinkOptionId", drink.id, { shouldValidate: true });
      form.setValue("abv", drink.abv * 100, { shouldValidate: true });
      form.setValue("sizeValue", drink.size_value, { shouldValidate: true });
      form.setValue("sizeUnit", drink.size_unit, { shouldValidate: true });
    }
  };

  async function onSubmit(data: DrinkLogFormValues) {
    if (!form.formState.isValid) {
      return;
    }

    setIsLoading(true);
    try {
      const newLog = {
        drink_option_id: data.drinkOptionId,
        logged_at: new Date().toISOString(),
      };
      await apiService.createDrinkLog(newLog.drink_option_id, newLog.logged_at);

      // addDrinkLog(newLog);
      form.reset();
      toast.success("Drink logged successfully");
      refreshDrinkLogs();
    } catch (error) {
      console.error("Failed to log drink:", error);
      toast.error("Failed to log drink");
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="drinkOptionId"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Select Drink</FormLabel>
              <Select
                onValueChange={onDrinkSelect}
                value={field.value?.toString() || ""}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a drink" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {drinkOptions.map((drink) => (
                    <SelectItem key={drink.id} value={drink.id.toString()}>
                      {drink.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="abv"
          render={({ field }) => (
            <FormItem>
              <FormLabel>ABV %</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  step="0.50"
                  placeholder="Alcohol percentage"
                  {...field}
                  value={field.value || ""}
                  onChange={(e) =>
                    field.onChange(
                      e.target.value ? parseFloat(e.target.value) : ""
                    )
                  }
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex gap-4">
          <FormField
            control={form.control}
            name="sizeValue"
            render={({ field }) => (
              <FormItem className="flex-1">
                <FormLabel>Size</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    placeholder="Amount"
                    {...field}
                    value={field.value || ""}
                    onChange={(e) =>
                      field.onChange(
                        e.target.value ? parseInt(e.target.value) : ""
                      )
                    }
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="sizeUnit"
            render={({ field }) => (
              <FormItem className="w-24">
                <FormLabel>Unit</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={isLoading || !form.formState.isValid}
        >
          {isLoading ? "Logging..." : "Log Drink"}
        </Button>
      </form>
    </Form>
  );
}
