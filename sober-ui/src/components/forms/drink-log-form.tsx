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

import { DrinkTemplate, DrinkLog } from "@/lib/types/api";
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

const createDrinkLogSchema = z.object({
  drinkTemplateId: z.string().min(1, "Please select a drink"),
  abv: z.number().min(0.01, "ABV must be greater than 0").max(100),
  sizeValue: z.number().min(1, "Size must be greater than 0"),
  sizeUnit: z.enum(["cl", "ml"], {
    errorMap: () => ({ message: "Please select a valid unit" }),
  }),
});

const editDrinkLogSchema = z.object({
  name: z.string().min(1, "Name is required"),
  abv: z.number().min(0.01, "ABV must be greater than 0").max(100),
  sizeValue: z.number().min(1, "Size must be greater than 0"),
  sizeUnit: z.enum(["cl", "ml"], {
    errorMap: () => ({ message: "Please select a valid unit" }),
  }),
});

type DrinkLogFormValues = z.infer<typeof createDrinkLogSchema> &
  z.infer<typeof editDrinkLogSchema>;

interface DrinkLogFormProps {
  initialDrinkLog?: DrinkLog;
  onCancel?: () => void;
  mode?: "create" | "edit";
}

export function DrinkLogForm({
  initialDrinkLog,
  onCancel,
  mode = "create",
}: DrinkLogFormProps) {
  const [drinkTemplates, setDrinkTemplates] = useState<DrinkTemplate[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<DrinkLogFormValues>({
    resolver: zodResolver(
      mode === "create" ? createDrinkLogSchema : editDrinkLogSchema
    ),
    defaultValues: {
      drinkTemplateId:
        mode === "create" ? initialDrinkLog?.id?.toString() : undefined,
      name: mode === "edit" ? initialDrinkLog?.name ?? "" : undefined,
      abv: initialDrinkLog ? initialDrinkLog.abv * 100 : 0,
      sizeValue: initialDrinkLog?.size_value ?? 0,
      sizeUnit: (initialDrinkLog?.size_unit as "cl" | "ml") ?? "cl",
    },
    mode: "onChange",
  });

  const { addDrinkLog, refreshDrinkLogs } = useDrinkLogs();

  useEffect(() => {
    apiService
      .getDrinkTemplates()
      .then((data) => {
        const templates = Array.isArray(data.drink_templates)
          ? data.drink_templates
          : [];
        setDrinkTemplates(templates);
      })
      .catch((error) => {
        console.error("Failed to fetch drink templates:", error);
        setDrinkTemplates([]);
      });
  }, []);

  const onDrinkSelect = (drinkId: string) => {
    const drink = drinkTemplates.find((d) => d.id === parseInt(drinkId));
    if (drink) {
      console.log({ drink });
      form.setValue("drinkTemplateId", drink.id.toString());
      form.setValue("abv", drink.abv * 100, { shouldValidate: true });
      form.setValue("sizeValue", drink.size_value, { shouldValidate: true });
      form.setValue("sizeUnit", drink.size_unit as "cl" | "ml", {
        shouldValidate: true,
      });
    }
  };

  async function onSubmit(data: DrinkLogFormValues) {
    if (!form.formState.isValid) return;

    setIsLoading(true);
    try {
      const drinkData = (() => {
        switch (mode) {
          case "create": {
            const selectedTemplate = drinkTemplates.find(
              (d) => d.id.toString() === data.drinkTemplateId
            );
            if (!selectedTemplate) {
              throw new Error("No drink template selected");
            }
            return {
              id: selectedTemplate.id,
              name: selectedTemplate.name,
              type: selectedTemplate.type,
              size_value: data.sizeValue,
              size_unit: data.sizeUnit,
              abv: data.abv / 100,
            };
          }
          case "edit": {
            if (!initialDrinkLog?.id) {
              throw new Error("No drink log ID provided for edit");
            }
            return {
              id: initialDrinkLog.id,
              name: data.name,
              type: initialDrinkLog?.type || "custom",
              size_value: data.sizeValue,
              size_unit: data.sizeUnit,
              abv: data.abv / 100,
            };
          }
        }
      })();

      if (mode === "edit" && initialDrinkLog) {
        await apiService.updateDrinkLog(drinkData);
        toast.success("Drink updated successfully");
      } else {
        await apiService.createDrinkLog(drinkData);
        toast.success("Drink logged successfully");
      }

      form.reset();
      refreshDrinkLogs();
      if (onCancel) onCancel();
    } catch (error) {
      console.error(`Failed to ${mode} drink:`, error);
      toast.error(`Failed to ${mode} drink`);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {mode === "create" ? (
          <FormField
            control={form.control}
            name="drinkTemplateId"
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
                    {drinkTemplates.map((drink) => (
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
        ) : (
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Drink Name</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        )}

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
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Unit" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="cl">cl</SelectItem>
                    <SelectItem value="ml">ml</SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <div className="flex gap-2">
          <Button
            type="submit"
            className="flex-1"
            disabled={isLoading || !form.formState.isValid}
          >
            {isLoading
              ? `${mode === "edit" ? "Updating..." : "Logging..."}`
              : `${mode === "edit" ? "Update" : "Log"} Drink`}
          </Button>
          {mode === "edit" && onCancel && (
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
          )}
        </div>
      </form>
    </Form>
  );
}
