create TABLE if not exists ratings (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL,
  shop_id bigint NOT NULL,
  ambiance_rating smallint,
  coffee_rating smallint,
  overall_rating smallint,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
