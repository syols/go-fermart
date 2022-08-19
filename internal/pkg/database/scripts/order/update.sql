INSERT INTO orders (order_id, user_id, accrual, order_status)
VALUES (:order_id, :user_id, :accrual, :order_status) ON CONFLICT (order_id) DO
    UPDATE
    SET user_id = excluded.user_id,
        accrual = excluded.accrual,
        order_status = excluded.order_status;