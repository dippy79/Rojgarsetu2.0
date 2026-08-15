package legal

const (
	// Disclaimer - IT Act 2000 Section 79 compliance
	Disclaimer = `
RojgarSetu is an intermediary platform under IT Act 2000 Section 79.
Job listings are sourced from official government portals (public records)
and authorized API partners. We do not host or modify original content.
For takedown requests: legal@rojgarsetu.in
Bot Policy: https://rojgarsetu.in/bot-policy
`

	// PrivacyPolicy - DPDP Act 2023 compliant statement
	PrivacyPolicy = `
RojgarSetu collects and processes personal data in accordance with
the Digital Personal Data Protection Act 2023. We process user data
solely for job matching and recruitment purposes. Users may request
data deletion or access by contacting privacy@rojgarsetu.in.
`

	// TermsOfService - basic terms
	TermsOfService = `
By using RojgarSetu, you agree to:
- Provide accurate information for job applications
- Not misuse the platform for fraudulent activities
- Respect intellectual property rights of job postings
- Accept that RojgarSetu acts as an intermediary only
`
)

// GetDisclaimer returns the IT Act compliance disclaimer
func GetDisclaimer() string {
	return Disclaimer
}

// GetPrivacyPolicy returns the DPDP Act 2023 privacy policy
func GetPrivacyPolicy() string {
	return PrivacyPolicy
}

// GetTermsOfService returns the terms of service
func GetTermsOfService() string {
	return TermsOfService
}
