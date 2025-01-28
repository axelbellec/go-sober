"use client";

import { useState, useEffect, useMemo } from "react";
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
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

import {
  DrinkTemplate,
  DrinkLog,
  ParseDrinkLogResponse,
} from "@/lib/types/api";
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
import { useDebounce } from "@/hooks/use-debounce";

const createDrinkLogSchema = z
  .object({
    drinkTemplateId: z.string().optional(),
    freeText: z.string().optional(),
    abv: z.number().min(0.01, "ABV must be greater than 0").max(100),
    sizeValue: z.number().min(1, "Size must be greater than 0"),
    sizeUnit: z.enum(["cl", "ml"], {
      errorMap: () => ({ message: "Please select a valid unit" }),
    }),
  })
  .refine((data) => data.drinkTemplateId || data.freeText, {
    message: "Please either select a drink or describe what you had",
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
  onDelete?: () => void;
}

const drinkLogPlaceholders = [
  "e.g., had a glass of red wine",
  "e.g., drank a pint of beer",
  "e.g., had a mojito cocktail",
  "e.g., enjoyed a gin and tonic",
  "e.g., had a shot of tequila",
  "e.g., drank a vodka soda",
  "e.g., had a glass of prosecco",
  "e.g., enjoyed an old fashioned",
  "e.g., had a margarita",
  "e.g., drank a whiskey neat",
];

const getRandomPlaceholder = () => {
  return drinkLogPlaceholders[
    Math.floor(Math.random() * drinkLogPlaceholders.length)
  ];
};

export function DrinkLogForm({
  initialDrinkLog,
  onCancel,
  mode = "create",
  onDelete,
}: DrinkLogFormProps) {
  const [drinkTemplates, setDrinkTemplates] = useState<DrinkTemplate[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isParsingDrink, setIsParsingDrink] = useState(false);
  const [placeholder, setPlaceholder] = useState("Describe your drink...");

  useEffect(() => {
    setPlaceholder(getRandomPlaceholder());
  }, []);

  const form = useForm<DrinkLogFormValues>({
    resolver: zodResolver(
      mode === "create" ? createDrinkLogSchema : editDrinkLogSchema
    ),
    defaultValues: {
      drinkTemplateId:
        mode === "create" ? initialDrinkLog?.id?.toString() : undefined,
      freeText: "",
      name: mode === "edit" ? initialDrinkLog?.name ?? "" : undefined,
      abv: initialDrinkLog ? initialDrinkLog.abv * 100 : 0,
      sizeValue: initialDrinkLog?.size_value ?? 0,
      sizeUnit: (initialDrinkLog?.size_unit as "cl" | "ml") ?? "cl",
    },
    mode: "onChange",
  });

  const { refreshDrinkLogs } = useDrinkLogs();

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

  const parseDrink = async (value: string) => {
    if (!value) return;

    setIsParsingDrink(true);
    try {
      const response: ParseDrinkLogResponse = await apiService.parseDrinkLog(
        value
      );

      form.setValue("drinkTemplateId", response.drink_template.id.toString(), {
        shouldValidate: true,
      });
      form.setValue("abv", response.drink_template.abv * 100, {
        shouldValidate: true,
      });
      form.setValue("sizeValue", response.drink_template.size_value, {
        shouldValidate: true,
      });
      form.setValue(
        "sizeUnit",
        response.drink_template.size_unit as "cl" | "ml",
        {
          shouldValidate: true,
        }
      );
    } catch (error) {
      console.error("Failed to parse drink:", error);
    } finally {
      setIsParsingDrink(false);
    }
  };

  const debouncedParseDrink = useDebounce(parseDrink, 400);

  const onFreeTextChange = (value: string) => {
    form.setValue("freeText", value);
    form.setValue("drinkTemplateId", undefined);
    debouncedParseDrink(value);
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
        toast.success(`Updated ${drinkData.name}`, {
          description: `Changed to ${drinkData.size_value}${
            drinkData.size_unit
          } at ${(drinkData.abv * 100).toFixed(1)}% ABV`,
          duration: 3000,
        });
      } else {
        await apiService.createDrinkLog(drinkData);
        toast.success(`Added ${drinkData.name} to your log`, {
          description: `${drinkData.size_value}${drinkData.size_unit} at ${(
            drinkData.abv * 100
          ).toFixed(1)}% ABV`,
          duration: 3000,
        });
      }

      form.reset();
      refreshDrinkLogs();
      if (onCancel) onCancel();
    } catch (error) {
      console.error(`Failed to ${mode} drink:`, error);
      toast.error(`Unable to ${mode === "edit" ? "update" : "log"} drink`, {
        description:
          "Something went wrong. Please check your input and try again.",
        duration: 5000,
      });
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {mode === "create" && (
          <>
            <FormField
              control={form.control}
              name="freeText"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Describe your drink</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <Input
                        placeholder={placeholder}
                        autoComplete="off"
                        {...field}
                        onChange={(e) => onFreeTextChange(e.target.value)}
                      />
                      {isParsingDrink && (
                        <div className="absolute right-3 top-1/2 -translate-y-1/2">
                          <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                        </div>
                      )}
                    </div>
                  </FormControl>
                  <FormDescription>
                    Describe what you had or select from templates below
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                  or select from templates
                </span>
              </div>
            </div>

            <FormField
              control={form.control}
              name="drinkTemplateId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Select Drink</FormLabel>
                  <Select
                    onValueChange={(value) => {
                      form.setValue("freeText", "");
                      onDrinkSelect(value);
                    }}
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
          </>
        )}

        {mode === "edit" && (
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
          {mode === "edit" && (
            <>
              {onCancel && (
                <Button type="button" variant="outline" onClick={onCancel}>
                  Cancel
                </Button>
              )}
              {onDelete && (
                <Button
                  type="button"
                  variant="destructive"
                  onClick={async () => {
                    if (initialDrinkLog) {
                      try {
                        toast.loading("Deleting drink...");
                        await apiService.deleteDrinkLog(initialDrinkLog.id);
                        toast.success(`Removed ${initialDrinkLog.name}`, {
                          description:
                            "The drink has been deleted from your log",
                          duration: 3000,
                        });
                        refreshDrinkLogs();
                        if (onDelete) onDelete();
                      } catch (error) {
                        toast.error("Failed to delete", {
                          description:
                            "Please try again or contact support if the problem persists",
                          duration: 5000,
                        });
                      }
                    }
                  }}
                >
                  Delete
                </Button>
              )}
            </>
          )}
        </div>
      </form>
    </Form>
  );
}
