# Receipt-Processor

## Overview
The Receipt-Processor is a simple service that processes receipts, calculates points based on specific rules, and provides information about receipts.

## Usage

### Start the server

- Clone the repository: https://github.com/william0228/Receipt-Processor.git
- Install Go dependencies
- run command: ```go run main.go``` (The server will run on http://localhost:8080)

## API Endpoints

#### 1. POST /receipts/process
This endpoint is used to submit a receipt for processing.

Example Request:

```text
GET http://localhost:8080/receipts/process
```

Request body example:

```json
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    }
  ],
  "total": "6.49"
}
```

Response (Success): The response will return the id of the processed receipt.

```json
{
  "id": "adb6b560-0eef-42bc-9d16-df48f30e89b2"
}
```


Response (Failure):

```json
{
  "error": "Invalid receipt"
}
```

#### 2. GET /receipts/{id}/points
This endpoint returns the points awarded for the receipt with the specified id.

Example Request:

```text
GET http://localhost:8080/receipts/{id}/points
```

Response:
```json
{
  "points": 100
}
```

Response (Not Found):

```json
{
  "error": "Receipt not found"
}
```

#### 3. GET /receipts/all
This endpoint returns a list of all receipts stored in the system.

Example Request:

```text
GET http://localhost:8080/receipts/all
```

Response:

```json
{
  "adb6b560-0eef-42bc-9d16-df48f30e89b2": {
    "retailer": "M&M Corner Market",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      {
        "shortDescription": "Mountain Dew 12PK",
        "price": "6.49"
      }
    ],
    "total": "6.49"
  }
}
```

## Points Calculation Rules

### The points awarded for a receipt are calculated based on the following rules:

- One point for every alphanumeric character in the retailer name.
- 50 points if the total amount on the receipt is a round dollar amount (i.e., has no cents).
- 25 points if the total amount is a multiple of 0.25.
- 5 points for every two items on the receipt.
- If the trimmed length of the item description is a multiple of 3: Multiply the price of the item by 0.2 and round up to the nearest integer. This is the number of points awarded for that item.
- 6 points if the day in the purchase date is odd (e.g., 1st, 3rd, 5th, etc.).
- 10 points if the purchase time is after 2:00 PM and before 4:00 PM.

### Example Calculation

```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```

```text
Total Points: 28
Breakdown:
     6 points - retailer name has 6 characters
    10 points - 5 items (2 pairs @ 5 points each)
     3 Points - "Emils Cheese Pizza" is 18 characters (a multiple of 3)
                item price of 12.25 * 0.2 = 2.45, rounded up is 3 points
     3 Points - "Klarbrunn 12-PK 12 FL OZ" is 24 characters (a multiple of 3)
                item price of 12.00 * 0.2 = 2.4, rounded up is 3 points
     6 points - purchase day is odd
  + ---------
  = 28 points
```

----

```json
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```

```text
Total Points: 109
Breakdown:
    50 points - total is a round dollar amount
    25 points - total is a multiple of 0.25
    14 points - retailer name (M&M Corner Market) has 14 alphanumeric characters
                note: '&' is not alphanumeric
    10 points - 2:33pm is between 2:00pm and 4:00pm
    10 points - 4 items (2 pairs @ 5 points each)
  + ---------
  = 109 points
```
