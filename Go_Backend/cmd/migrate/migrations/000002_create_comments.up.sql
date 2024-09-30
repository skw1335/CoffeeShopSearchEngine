CREATE TABLE if not exists comments (
  id bigserial PRIMARY KEY,
  content text NOT NULL,
  title text NOT NULL,
  user_id bigint NOT NULL,
  shop_id bigint NOT NULL,
  tags       text[],
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
