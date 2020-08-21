CREATE TABLE IF NOT EXISTS rewards (
    blk_lvl UInt64,
    blk_epoch UInt64,
    created_at DateTime,
    rwd_amount UInt64,
    reg_entity_address  FixedString(46)
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(created_at)
ORDER BY (reg_entity_address, blk_lvl);

CREATE VIEW validator_rewards_stat_view AS
select *
from (
       select *
       from (select reg_entity_address, sum(rwd_amount) total_amount
             from oasis.rewards
             group by reg_entity_address) total
              ANY
              LEFT JOIN (select reg_entity_address, sum(rwd_amount) day_amount
                         from oasis.rewards
                         where created_at >= toStartOfDay(now())
                         group by reg_entity_address) USING reg_entity_address
       )
       ANY
       LEFT JOIN (
  select *
  from (
         select reg_entity_address, sum(rwd_amount) week_amount
         from oasis.rewards
         where created_at >= toStartOfWeek(now())
         group by reg_entity_address) week
         ANY
         LEFT JOIN (select reg_entity_address, sum(rwd_amount) month_amount
                    from oasis.rewards
                    where created_at >= toStartOfMonth(now())
                    group by reg_entity_address) USING reg_entity_address
  ) USING reg_entity_address;