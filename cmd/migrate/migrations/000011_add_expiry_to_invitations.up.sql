ALTER TABLE user_invitations
ADD COLUMN expiry timestamp(0)
WITH
    TIME ZONE NOT NULL;