SELECT id, user_id, number, score, status, action, ctime FROM orders
WHERE number = :number;
