CREATE TABLE IF NOT EXISTS public.clientes (
  id SERIAL PRIMARY KEY,
  limite INT NOT NULL,
  saldo INT NOT NULL
);

CREATE TYPE tipo_transacao AS ENUM ('c', 'd');

CREATE TABLE IF NOT EXISTS public.transacoes (
  id SERIAL PRIMARY KEY,
  cliente_id INT REFERENCES clientes(id) NOT NULL,
  tipo tipo_transacao NOT NULL,
  valor INT NOT NULL,
  descricao VARCHAR(10) NOT NULL,
  realizada_em timestamp with time zone
);

INSERT INTO public.clientes (limite, saldo)
VALUES
  (100000, 0),
  (80000, 0),
  (1000000, 0),
  (10000000, 0),
  (500000, 0);