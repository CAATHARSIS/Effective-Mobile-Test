CREATE TABLE IF NOT EXISTS subscription_record (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_name TEXT NOT NULL,
    price_per_mounth INT NOT NULL CHECK(price_per_mounth >= 0),
    user_uuid UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    CHECK(end_date IS NULL OR end_date >= start_date)
);