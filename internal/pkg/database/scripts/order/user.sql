SELECT order_id, user_id, accrual, order_status, uploaded FROM orders
WHERE user_id = :user_id ORDER BY uploaded;