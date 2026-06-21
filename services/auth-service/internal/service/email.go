package service

// SendVerificationEmail sends an email verification link to the user (stub — requires SendGrid configuration)
func SendVerificationEmail(email, token string) error {
	// TODO: integrate SendGrid API
	return nil
}

// SendPasswordResetEmail sends a password reset link to the user (stub — requires SendGrid configuration)
func SendPasswordResetEmail(email, token string) error {
	// TODO: integrate SendGrid API
	return nil
}
