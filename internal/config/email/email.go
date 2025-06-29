package email

import "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"

func PrepareEmailData(cfg *config.Config) map[string]config.EmailData {
	emailDataMap := make(map[string]config.EmailData)

	for _, email := range cfg.Emails {
		memail := email
		memail.From = cfg.EmailService.Sender
		emailDataMap[memail.Name] = memail
	}

	return emailDataMap
}
