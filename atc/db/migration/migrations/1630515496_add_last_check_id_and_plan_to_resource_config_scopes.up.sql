ALTER TABLE resource_config_scopes
    ADD COLUMN last_check_build_id bigint,
    ADD COLUMN last_check_build_plan json NULL DEFAULT '{}'::json;

-- This index should not be unique, because a build may check multiple resources and resource types.
CREATE INDEX resource_config_scopes_last_check_build_id_idx ON resource_config_scopes (last_check_build_id);
