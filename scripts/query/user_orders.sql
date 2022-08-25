SELECT user_id, number, score, status, ctime FROM orders
WHERE user_id = :user_id AND action = :action ORDER BY ctime;
