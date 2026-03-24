-- Create birthday_groups table
CREATE TABLE IF NOT EXISTS birthday_groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    birthday_user_id BIGINT NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL DEFAULT 'upcoming',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create birthday_group_members table (M2M)
CREATE TABLE IF NOT EXISTS birthday_group_members (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES birthday_groups(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(group_id, user_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_birthday_groups_birthday_user_id ON birthday_groups(birthday_user_id);
CREATE INDEX IF NOT EXISTS idx_birthday_groups_status ON birthday_groups(status);
CREATE INDEX IF NOT EXISTS idx_birthday_group_members_group_id ON birthday_group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_birthday_group_members_user_id ON birthday_group_members(user_id);
