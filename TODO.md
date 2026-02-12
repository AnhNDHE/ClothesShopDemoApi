# Email Verification Implementation

## Completed Tasks
- [x] Update migrations/000001_init_tables.up.sql: Add base fields to users table, set is_active DEFAULT false
- [x] Update internal/config/config.go: Add SMTP config fields and load them in LoadConfig
- [x] Create internal/services/email_service.go: Email service using net/smtp to send verification emails
- [x] Update internal/repositories/user_repository.go: Modify CreateUser to set is_active=false, update GetUserByEmail to include is_active, add ActivateUser method
- [x] Update internal/handlers/auth_handler.go: Modify Register to send verification email instead of activating, add VerifyEmail handler, update Login to check is_active
- [x] Update internal/routes/routes.go: Add /auth/verify-email route
- [x] Update Swagger docs for new endpoint

## Pending Tasks
- [x] Run migration to update database
- [x] Test email sending and verification flow
