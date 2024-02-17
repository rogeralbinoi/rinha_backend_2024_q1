BEGIN;

CREATE TABLE IF NOT EXISTS public.clientes (
  id SERIAL PRIMARY KEY,
  limite INT NOT NULL,
  saldo_inicial INT NOT NULL
);

INSERT INTO public.clientes (limite, saldo_inicial)
VALUES
  (100000, 0),
  (80000, 0),
  (1000000, 0),
  (10000000, 0),
  (500000, 0);

SELECT * FROM public.clientes;

COMMIT;