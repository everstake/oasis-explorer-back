CREATE TABLE IF NOT EXISTS tasks (
   tsk_id             serial constraint tasks_pk primary key,
   tsk_active         boolean default TRUE,
   tsk_title          VARCHAR(255) NOT NULL,
   tsk_start_height   int        NOT NULL,
   tsk_current_height int      NOT NULL,
   tsk_end_height     int      NOT NULL,
   tsk_batch          int      NOT NULL
);