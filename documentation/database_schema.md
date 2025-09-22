                                       Table "public.files"
    Column    |          Type          | Collation | Nullable |              Default
--------------+------------------------+-----------+----------+-----------------------------------
 id           | integer                |           | not null | nextval('files_id_seq'::regclass)
 filename     | character varying(255) |           | not null |
 hash         | character varying(255) |           | not null |
 storage_path | text                   |           | not null |
 size         | bigint                 |           |          |
 mime_type    | character varying(255) |           | not null |
 ref_count    | integer                |           |          | 1


 fvs_mgmt_sys=> \d users;
                                         Table "public.users"
      Column      |          Type          | Collation | Nullable |              Default
------------------+------------------------+-----------+----------+-----------------------------------
 id               | integer                |           | not null | nextval('users_id_seq'::regclass)
 username         | character varying(50)  |           | not null |
 email            | character varying(100) |           | not null |
 password         | text                   |           | not null |
 actual_storage   | bigint                 |           | not null | 0
 expected_storage | bigint                 |           | not null | 0


                                          Table "public.user_files"
     Column     |           Type           | Collation | Nullable |                Default
----------------+--------------------------+-----------+----------+----------------------------------------
 id             | integer                  |           | not null | nextval('user_files_id_seq'::regclass)
 user_id        | integer                  |           | not null |
 file_id        | integer                  |           | not null |
 file_name      | character varying(255)   |           |          |
 uploaded_at    | timestamp with time zone |           |          |
 download_times | integer                  |           |          |
 is_owner       | boolean                  |           |          | false
 visibility     | text                     |           |          | 'private'::text
 public_token   | text                     |           |          |