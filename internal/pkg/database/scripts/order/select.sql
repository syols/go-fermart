SELECT order_id, user_id, accrual, order_status, uploaded FROM orders
WHERE order_id = :order_id;