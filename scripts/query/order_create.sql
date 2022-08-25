INSERT INTO orders (user_id, number, score, status, action)
VALUES (:user_id, :number, :score, :status, :action) ON CONFLICT (number) DO
    UPDATE
    SET user_id = excluded.user_id,
        number = excluded.number,
        score = excluded.score,
        status = excluded.status,
        action = excluded.action;