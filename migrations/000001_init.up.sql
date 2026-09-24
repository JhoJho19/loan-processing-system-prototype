BEGIN;

CREATE TYPE employment_type AS ENUM (
    'EMPLOYED',
    'SELF_EMPLOYED',
    'UNEMPLOYED',
    'STUDENT',
    'RETIRED'
);

CREATE TYPE application_status AS ENUM (
    'NEW',
    'VALIDATING',
    'SCORING',
    'APPROVED',
    'REJECTED'
);

CREATE TABLE customers (
    id UUID
        CONSTRAINT pk_customers PRIMARY KEY,

    first_name VARCHAR(100)
        CONSTRAINT nn_customers_first_name NOT NULL,

    last_name VARCHAR(100)
        CONSTRAINT nn_customers_last_name NOT NULL,

    birth_date DATE
        CONSTRAINT nn_customers_birth_date NOT NULL,

    monthly_income INTEGER
        CONSTRAINT nn_customers_monthly_income NOT NULL
        CONSTRAINT chk_customers_monthly_income_positive
            CHECK (monthly_income > 0),

    employment_type employment_type
        CONSTRAINT nn_customers_employment_type NOT NULL,

    email VARCHAR(255)
        CONSTRAINT nn_customers_email NOT NULL
        CONSTRAINT chk_customers_email_format
            CHECK (email LIKE '%_@_%._%'),

    phone VARCHAR(30)
        CONSTRAINT nn_customers_phone NOT NULL
        CONSTRAINT chk_customers_phone_format
            CHECK (phone ~ '^\+[1-9][0-9]{7,14}$'),

    address TEXT
        CONSTRAINT nn_customers_address NOT NULL
);

CREATE TABLE loan_applications (
    id UUID
        CONSTRAINT pk_loan_applications PRIMARY KEY,

    customer_id UUID
        CONSTRAINT nn_loan_applications_customer_id NOT NULL
        CONSTRAINT fk_loan_applications_customer_id
            REFERENCES customers(id),

    requested_amount INTEGER
        CONSTRAINT nn_loan_applications_requested_amount NOT NULL
        CONSTRAINT chk_loan_applications_requested_amount_positive
            CHECK (requested_amount > 0),

    term_months INTEGER
        CONSTRAINT nn_loan_applications_term_months NOT NULL
        CONSTRAINT chk_loan_applications_term_months_positive
            CHECK (term_months > 0),

    status application_status
        CONSTRAINT nn_loan_applications_status NOT NULL
        DEFAULT 'NEW',

    created_at TIMESTAMPTZ
        CONSTRAINT nn_loan_applications_created_at NOT NULL
        DEFAULT now(),

    updated_at TIMESTAMPTZ
        CONSTRAINT nn_loan_applications_updated_at NOT NULL
        DEFAULT now()
);

CREATE INDEX idx_loan_applications_customer_id
    ON loan_applications(customer_id);

COMMIT;