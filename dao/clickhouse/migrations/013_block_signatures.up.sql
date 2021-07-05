-- PARTITIONED by 500000
CREATE MATERIALIZED VIEW IF NOT EXISTS block_signatures_count_mv
  ENGINE = SummingMergeTree()
    PARTITION BY intDiv(blk_lvl,500000)
    ORDER BY (blk_lvl) POPULATE AS
select blk_lvl,
       count()                     signatures
from block_signatures
group by blk_lvl;

CREATE VIEW IF NOT EXISTS block_signatures_count_view AS
select blk_lvl,
       sum(signatures) sig_count
from block_signatures_count_mv
group by blk_lvl;

DROP TABLE blocks_sig_count;