SELECT orders.user_id,
       SUM(CASE action WHEN 'PURCHASE' THEN score ELSE 0 END) -
       SUM(CASE action WHEN 'WITHDRAW' THEN score ELSE 0 END) AS current,
       SUM(CASE action WHEN 'WITHDRAW' THEN score ELSE 0 END) AS withdrawn
FROM orders WHERE user_id = :user_id AND status = 'PROCESSED' GROUP BY user_id;