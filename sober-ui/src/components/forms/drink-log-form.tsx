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

const drinkLogSchema = z.object({
  drinkOptionId: z.number(),
  abv: z.number().min(0).max(100),
  sizeValue: z.number().min(0),
  sizeUnit: z.string(),
});

type DrinkLogFormValues = z.infer<typeof drinkLogSchema>;

export function DrinkLogForm() {
  const [drinkOptions, setDrinkOptions] = useState<DrinkOption[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<DrinkLogFormValues>({
    resolver: zodResolver(drinkLogSchema),
    defaultValues: {
      drinkOptionId: 1,
      abv: 0,
      sizeValue: 0,
      sizeUnit: "cl",
    },
  });

  useEffect(() => {
    fetch("http://localhost:3000/api/v1/drink-options")
      .then((res) => res.json())
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
      form.setValue("drinkOptionId", drink.id);
      form.setValue("abv", drink.abv * 100);
      form.setValue("sizeValue", drink.size_value);
      form.setValue("sizeUnit", drink.size_unit);
    }
  };

  async function onSubmit(data: DrinkLogFormValues) {
    setIsLoading(true);
    try {
      const response = await fetch("http://localhost:3000/api/v1/drink-logs", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({
          drink_option_id: data.drinkOptionId,
          logged_at: new Date().toISOString(),
        }),
      });

      if (response.ok) {
        // Reset form and show success message
        form.reset();
      }
    } catch (error) {
      console.error("Failed to log drink:", error);
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
                value={field.value.toString()}
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

        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Logging..." : "Log Drink"}
        </Button>
      </form>
    </Form>
  );
}
