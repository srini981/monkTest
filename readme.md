# Handler Scenarios Documentation

This document outlines all scenarios for each handler based on the requirements provided in the task document. Each handler is detailed with its purpose, inputs, outputs, and all possible test cases.

---

## **1. `POST /coupons` - Create a New Coupon**

### **Purpose**
Create a new coupon with specified conditions, type, and expiration date.

### **Inputs**
- **Name**: Name of the coupon
- **Type**: Specifies the type of coupon (cart-wise, product-wise, bxgy).
- **Details**: Contains type-specific parameters like discount percentage, thresholds, product IDs, etc.
- **ExpiresAt**: Timestamp indicating when the coupon expires.

### **Scenarios**
1. **Valid Coupon Creation**
   - Inputs: Valid type, details, and future expiration date.
   - Expected Outcome: Coupon created successfully.

2. **Invalid Coupon Type**
   - Inputs: Unsupported type (e.g., "unknown").
   - Expected Outcome: Error response with validation message.

3. **Expired Coupon on Creation**
   - Inputs: Expiration date in the past.
   - Expected Outcome: Error response stating "coupon has expired."

4. **Missing Required Fields**
   - Inputs: Missing type, details, or expiration date.
   - Expected Outcome: Error response indicating invalid request payload.

5. **Malformed Input Data**
   - Inputs: Incorrect JSON structure or invalid data types.
   - Expected Outcome: Error response indicating invalid payload.

6. **Duplicate token creation**
   - Inputs: two asyn calls hits with the same JSON data cause duplicate tokens.
   - Expected Outcome: Error response indicating duplicate coupon.

. **Database Connectivity Issues**
   - Inputs: Simulate database errors (e.g., MongoDB server down).
   - Expected Outcome: Error response indicating failure to create coupons.

---

## **2. `GET /coupons` - Retrieve All Coupons**

### **Purpose**
Fetch all coupons from the database.

### **Scenarios**
1. **Successful Retrieval**
   - Inputs: No additional parameters.
   - Expected Outcome: List of all coupons, including active and expired ones.

2. **No Coupons in Database**
   - Inputs: Empty database.
   - Expected Outcome: Empty list response.

3. **Database Connectivity Issues**
   - Inputs: Simulate database errors (e.g., MongoDB server down).
   - Expected Outcome: Error response indicating failure to retrieve coupons.

---

## **3. `GET /coupons/{id}` - Retrieve a Specific Coupon**

### **Purpose**
Fetch a coupon by its unique ID.

### **Scenarios**
1. **Valid Coupon ID**
   - Inputs: Existing coupon ID.
   - Expected Outcome: Return coupon details.

2. **Non-Existent Coupon ID**
   - Inputs: Validly formatted but non-existent coupon ID.
   - Expected Outcome: Error response stating "coupon not found."

3. **Invalid Coupon ID Format**
   - Inputs: Malformed or incorrect ID format (e.g., not an ObjectID).
   - Expected Outcome: Error response indicating invalid ID.

4. **Database Connectivity Issues**
   - Inputs: Simulate database errors (e.g., MongoDB server down).
   - Expected Outcome: Error response indicating failure to retrieve coupons.

---

## **4. `PUT /coupons/{id}` - Update a Coupon**

### **Purpose**
Update the details or expiration date of an existing coupon.

### **Scenarios**
1. **Valid Update Request**
   - Inputs: Existing coupon ID with valid updates.
   - Expected Outcome: Coupon updated successfully.

2. **Non-Existent Coupon ID**
   - Inputs: Validly formatted but non-existent coupon ID.
   - Expected Outcome: Error response stating "coupon not found."

3. **Invalid Input Data**
   - Inputs: Incorrect JSON structure or invalid data types.
   - Expected Outcome: Error response indicating invalid payload.

4. **Missing Required Fields**
   - Inputs: Missing type, details, or expiration date.
   - Expected Outcome: Error response indicating invalid request payload.

4. **Invalid Coupon ID Format**
   - Inputs: Malformed or incorrect ID format.
   - Expected Outcome: Error response indicating invalid ID.

5. **Database Connectivity Issues**
   - Inputs: Simulate database errors (e.g., MongoDB server down).
   - Expected Outcome: Error response indicating failure to retrieve coupons.

---

## **5. `DELETE /coupons/{id}` - Delete a Coupon**

### **Purpose**
Delete a coupon by its unique ID.

### **Scenarios**
1. **Valid Deletion Request**
   - Inputs: Existing coupon ID.
   - Expected Outcome: Coupon deleted successfully.

2. **Non-Existent Coupon ID**
   - Inputs: Validly formatted but non-existent coupon ID.
   - Expected Outcome: Error response stating "coupon not found."

3. **Invalid Coupon ID Format**
   - Inputs: Malformed or incorrect ID format.
   - Expected Outcome: Error response indicating invalid ID.

4. **Database Connectivity Issues**
   - Inputs: Simulate database errors (e.g., MongoDB server down).
   - Expected Outcome: Error response indicating failure to retrieve coupons.

---

## **6. `POST /applicable-coupons` - Fetch Applicable Coupons**

### **Purpose**
Retrieve all coupons applicable to a given cart and calculate the total discount for each.

### **Scenarios**
1. **Applicable Cart-Wise Coupon**
   - Inputs: Cart total exceeding coupon threshold.
   - Expected Outcome: Coupon listed with calculated discount.

2. **Applicable Product-Wise Coupon**
   - Inputs: Cart containing eligible products.
   - Expected Outcome: Coupon listed with calculated discount.

3. **Applicable BxGy Coupon**
   - Inputs: Cart meeting "buy" requirements for free products.
   - Expected Outcome: Coupon listed with calculated discount for free products.

4. **No Applicable Coupons**
   - Inputs: Cart not meeting any coupon conditions.
   - Expected Outcome: Empty applicable coupons list.

5. **Expired Coupons**
   - Inputs: Cart with eligible but expired coupons.
   - Expected Outcome: Expired coupons are excluded from the list.

6. **Database Connectivity Issues**
   - Inputs: Simulate database errors (e.g., MongoDB server down).
   - Expected Outcome: Error response indicating failure to retrieve coupons.

7. **Invalid Input Data**
   - Inputs: Incorrect JSON structure or invalid data types.
   - Expected Outcome: Error response indicating invalid payload.


Yet to be implemented

8. **coupon expiry after its applied in cart**
   - details: after the coupon is applied on a cart ,when reloading the cart , check if there are any expired coupon in the cart

9. **coupon expired in cart when its applied**
   - details: when applying the coupon on the cart check if it expiry within a 30 secs to ensure we dont apply the coupon on the cart / and also check the coupons once when reloading the cart/ checking out the cart for processing the order
   
10. **invalid due to a product out of stock**
   - details: after a coupon is applied in a cart , when reloading the cart, check whether any products is out of stock and the coupon becomes invalid

---

## **7. `POST /apply-coupon/{id}` - Apply a Specific Coupon**

### **Purpose**
Apply a specified coupon to a cart and return the updated cart with discounts applied.

### **Scenarios**
1. **Valid Cart-Wise Coupon Application**
   - Inputs: Cart total exceeding coupon threshold.
   - Expected Outcome: Updated cart with discounted total.

2. **Valid Product-Wise Coupon Application**
   - Inputs: Cart containing eligible products.
   - Expected Outcome: Updated cart with discounted product prices.

3. **Valid BxGy Coupon Application**
   - Inputs: Cart meeting "buy" requirements for free products.
   - Expected Outcome: Updated cart with free products added.

4. **Expired Coupon**
   - Inputs: Validly formatted but expired coupon.
   - Expected Outcome: Error response stating "coupon has expired."

5. **Invalid Coupon ID**
   - Inputs: Non-existent or malformed coupon ID.
   - Expected Outcome: Error response indicating invalid ID.

6. **Insufficient Cart Conditions**
   - Inputs: Cart not meeting coupon conditions (e.g., insufficient total or product quantity).
   - Expected Outcome: Error response stating "conditions not met."

7. **Database Connectivity Issues**
   - Inputs: Simulate database errors (e.g., MongoDB server down).
   - Expected Outcome: Error response indicating failure to retrieve coupons.

8. **Invalid Input Data**
   - Inputs: Incorrect JSON structure or invalid data types.
   - Expected Outcome: Error response indicating invalid payload.

Yet to be implemented

9. **coupon expiry after its applied in cart**
   - details: after the coupon is applied on a cart ,when reloading the cart , check if there are any expired coupon in the cart

10. **coupon expiredin cart when its applied**
   - details: when applying the coupon on the cart check if it expiry within a 30 secs to ensure we dont apply the coupon on the cart / and also check the coupons once when reloading the cart/ checking out the cart for processing the order
   
11. **invalid due to a product out of stock**
   - details: after a coupon is applied in a cart , when reloading the cart, check whether any products is out of stock and the coupon becomes invalid

---

This document serves as a comprehensive reference for all scenarios  tested for the handlers. Adjustments can be made based on further requirements or edge cases.


