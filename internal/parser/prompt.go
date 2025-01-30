package parser

const BEVERAGE_PARSER_PROMPT = `
# Instruction

Given a text description of a beverage, analyze the content and extract all relevant beverage information. Format the response as a JSON object containing an array of beverages with their complete details. Ensure all mandatory fields are populated and optional fields are included when the information is available in the input text.


# Requirements

Your response should respect the following requirements:
1) The model must correctly identify the type of beverage (beer, wine, champagne, etc.) from the input text. This is a fundamental requirement as it affects the entire parsing logic.
2) The model must accurately extract and standardize container volumes from the text, handling different units (ml, cl, L) and converting them appropriately. The volume must be split into value and unit components.
3) The model must detect and extract alcohol content when present, handling both percentage (%) and degree (°) symbols. The output should be standardized to percentage format. If the alcohol content is not explicitly stated, it should be set to -1.
4) The model must identify the type of container (bottle, glass, can) when mentioned in the text. This field is optional in the output but should be accurate when provided. 
5) The model must determine the number of beverages mentioned in the text. Default to 1 if not explicitly stated.
6) The model must generate a standardized name that includes the beverage name, volume, and alcohol content in a consistent format (e.g., "Heineken Beer, 330ml, 5%").

# Example Outputs

# Instruction

Given a text description of a beverage, analyze the content and extract all relevant beverage information. Format the response as a JSON object containing an array of beverages with their complete details. Ensure all mandatory fields are populated and optional fields are included when the information is available in the input text.


# Requirements

Your response should respect the following requirements:
1) The model must correctly identify the type of beverage (beer, wine, champagne, etc.) from the input text. This is a fundamental requirement as it affects the entire parsing logic.
2) The model must accurately extract and standardize container volumes from the text, handling different units (ml, cl, L) and converting them appropriately. The volume must be split into value and unit components.
3) The model must detect and extract alcohol content when present, handling both percentage (%) and degree (°) symbols. The output should be standardized to percentage format. If the alcohol content is not explicitly stated, it should be set to -1.
4) The model must identify the type of container (bottle, glass, can) when mentioned in the text. This field is optional in the output but should be accurate when provided. 
5) The model must determine the number of beverages mentioned in the text. Default to 1 if not explicitly stated.
6) The model must generate a standardized name that includes the beverage name, volume, and alcohol content in a consistent format (e.g., "Heineken Beer, 330ml, 5%").

# Example Outputs

## Example 1

Input: {"text": "pint of stout 4.2%"}
Output: {"beverages": [{"name": "Stout, 500ml, 4.2%", "container_volume_value": "500", "container_volume_unit": "ml", "alcohol_content": "4.2%", "quantity": 1, "type": "beer"}]}

## Example 2

Input: {"text": "house special IPA on tap 400ml 6.2%"}
Output: {"beverages": [{"name": "IPA, 400ml, 6.2%", "container_volume_value": "400", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "6.2%", "quantity": 1, "type": "beer"}]}

## Example 3

Input: {"text": "i had this mojito cocktail once again"}
Output: {"beverages": [{"name": "Mojito Cocktail, 240ml, 10%", "container_volume_value": "240", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "10%", "quantity": 1, "type": "cocktail"}]}

## Example 4

Input: {"text": "Ordered a glass of sparkling wine, 200ml, 11%."}
Output: {"beverages": [{"name": "Sparkling Wine, 200ml, 11%", "container_volume_value": "200", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "11%", "quantity": 1, "type": "wine"}]}

## Example 5

Input: {"text": "2 cans of pale ale 440ml each, 5.5%"}
Output: {"beverages": [{"name": "Pale Ale, 440ml, 5.5%", "container_volume_value": "440", "container_volume_unit": "ml", "container_type": "can", "alcohol_content": "5.5%", "quantity": 2, "type": "beer"}]}

## Example 6

Input: {"text": "Sipped on a 375ml bottle of Riesling, 12% alcohol."}
Output: {"beverages": [{"name": "Riesling Wine, 375ml, 12%", "container_volume_value": "375", "container_volume_unit": "ml", "container_type": "bottle", "alcohol_content": "12%", "quantity": 1, "type": "wine"}]}

## Example 7

Input: {"text": "one can of cider 500ml"}
Output: {"beverages": [{"name": "Cider, 500ml", "container_volume_value": "500", "container_volume_unit": "ml", "container_type": "can", "alcohol_content": "-1", "quantity": 1, "type": "cider"}]}

## Example 8

Input: {"text": "half a liter of Belgian ale, 500ml at 7%"}
Output: {"beverages": [{"name": "Belgian Ale, 500ml, 7%", "container_volume_value": "500", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "7%", "quantity": 1, "type": "beer"}]}

## Example 9

Input: {"text": "had a shot of tequila 25ml 38%"}
Output: {"beverages": [{"name": "Tequila, 25ml, 38%", "container_volume_value": "25", "container_volume_unit": "ml", "container_type": "shot", "alcohol_content": "38%", "quantity": 1, "type": "spirit"}]}

## Example 10

Input: {"text": "Drank a glass of Merlot, 150ml, 14%."}
Output: {"beverages": [{"name": "Merlot Wine, 150ml, 14%", "container_volume_value": "150", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "14%", "quantity": 1, "type": "wine"}]}

## Example 11

Input: {"text": "I sipped on a craft IPA, 473ml, 6.5%."}
Output: {"beverages": [{"name": "IPA, 473ml, 6.5%", "container_volume_value": "473", "container_volume_unit": "ml", "container_type": "can", "alcohol_content": "6.5%", "quantity": 1, "type": "beer"}]}

## Example 12

Input: {"text": "one IPA 33cl 6.7\u00b0"}
Output: {"beverages": [{"name": "IPA Beer, 33cl, 6.7%", "container_volume_value": "33", "container_volume_unit": "cl", "alcohol_content": "6.7%", "quantity": 1, "type": "beer"}]}

## Example 13

Input: {"text": "half pint of local IPA 6.5%"}
Output: {"beverages": [{"name": "Local IPA, 250ml, 6.5%", "container_volume_value": "250", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "6.5%", "quantity": 1, "type": "beer"}]}

## Example 14

Input: {"text": "I got a pint of Guinness, 568ml, 4.2%."}
Output: {"beverages": [{"name": "Guinness, 568ml, 4.2%", "container_volume_value": "568", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "4.2%", "quantity": 1, "type": "beer"}]}

## Example 15

Input: {"text": "one IPA 33cl 6.7\u00b0"}
Output: {"beverages": [{"name": "IPA Beer, 330ml, 6.7%", "container_volume_value": "330", "container_volume_unit": "ml", "container_type": "bottle", "alcohol_content": "6.7%", "quantity": 1, "type": "beer"}]}

## Example 16

Input: {"text": "Had a bottle of Chardonnay, 750ml, 13%."}
Output: {"beverages": [{"name": "Chardonnay Wine, 750ml, 13%", "container_volume_value": "750", "container_volume_unit": "ml", "container_type": "bottle", "alcohol_content": "13%", "quantity": 1, "type": "wine"}]}

## Example 17

Input: {"text": "one IPA 33cl 6.7\u00b0"}
Output: {"beverages": [{"name": "IPA, 330ml, 6.7%", "container_volume_value": "330", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "6.7%", "quantity": 1, "type": "beer"}]}

## Example 18

Input: {"text": "I had a nice Coors Light, 355ml, 4.2%."}
Output: {"beverages": [{"name": "Coors Light, 355ml, 4.2%", "container_volume_value": "355", "container_volume_unit": "ml", "container_type": "can", "alcohol_content": "4.2%", "quantity": 1, "type": "beer"}]}

## Example 19

Input: {"text": "I enjoyed a Stella Artois, 300ml, 5%."}
Output: {"beverages": [{"name": "Stella Artois, 300ml, 5%", "container_volume_value": "300", "container_volume_unit": "ml", "container_type": "bottle", "alcohol_content": "5%", "quantity": 1, "type": "beer"}]}

## Example 20

Input: {"text": "verre de ros\u00e9 ch\u00e2teau margaux 12\u00b0 15cl"}
Output: {"beverages": [{"name": "Ch\u00e2teau Margaux Ros\u00e9, 15cl, 12%", "container_volume_value": "150", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "12%", "quantity": 1, "type": "wine"}]}

## Example 21

Input: {"text": "draft lager, 500ml at 5.0%"}
Output: {"beverages": [{"name": "Draft Lager, 500ml, 5%", "container_volume_value": "500", "container_volume_unit": "ml", "alcohol_content": "5%", "quantity": 1, "type": "beer"}]}

## Example 22

Input: {"text": "2 cans of pale ale 440ml each, 5.5%"}
Output: {"beverages": [{"name": "Pale Ale, 440ml, 5.5%", "container_volume_value": "440", "container_volume_unit": "ml", "container_type": "can", "alcohol_content": "5.5%", "quantity": 2, "type": "beer"}]}

## Example 23

Input: {"text": "glass of red wine 175ml"}
Output: {"beverages": [{"name": "Red Wine, 175ml", "container_volume_value": "175", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "-1", "quantity": 1, "type": "wine"}]}

## Example 24

Input: {"text": "large glass of Merlot 250ml 13.5%"}
Output: {"beverages": [{"name": "Merlot Wine, 250ml, 13.5%", "container_volume_value": "250", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "13.5%", "quantity": 1, "type": "wine"}]}

## Example 25

Input: {"text": "bottle of hefeweizen 500ml 5.4% in a pint glass"}
Output: {"beverages": [{"name": "Hefeweizen Beer, 500ml, 5.4%", "container_volume_value": "500", "container_volume_unit": "ml", "container_type": "glass", "alcohol_content": "5.4%", "quantity": 1, "type": "beer"}]}

`
