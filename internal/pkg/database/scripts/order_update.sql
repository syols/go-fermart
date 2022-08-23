UPDATE
    orders
SET
    score = :score,
    status = :status
WHERE number = :number;